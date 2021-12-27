// Links a Chat-API instance to a database.

package main

import (
	"context"
	"log"
	"time"
)

type ChatAPIDB struct {
	ChatAPI *ChatAPI
	DB      *Database
	// Interval between forced updates.
	ReadMessageInterval time.Duration

	// Last message number in the database.
	// Used by CopyNewMessages to know where to start.
	LastMessageNumber int64
	// Send to this channel to force an update.
	UpdateC chan struct{}
	// The capacity is how many forced updates may happen without Active being called.
	ActiveC chan struct{}
}

func NewChatAPIDB(chatAPI *ChatAPI, db *Database) *ChatAPIDB {
	return &ChatAPIDB{
		ChatAPI:             chatAPI,
		DB:                  db,
		ReadMessageInterval: 10 * time.Second,
		ActiveC:             make(chan struct{}, 6),
		UpdateC:             make(chan struct{}, 1),
	}
}

// Start initializes an instance of ChatAPIDB.
func (wa *ChatAPIDB) Start(ctx context.Context) error {
	// Run one update on start.
	wa.UpdateC <- struct{}{}

	// Get last message number.
	number, err := wa.DB.GetLastMessageNumber(ctx)
	if err != nil {
		return err
	}

	// Get chats.
	err = wa.CopyChats(ctx)
	if err != nil {
		return err
	}

	// Start by copying the last 10 messages.
	if number > 10 {
		wa.LastMessageNumber = number - 10
	}

	// Start background goroutines.
	go wa.CopyNewMessagesLoop(ctx)
	go wa.CopyChatsLoop(ctx)

	// Started.
	return nil
}

// Active informs that automatic forced updates should continue.
func (wa *ChatAPIDB) Active() {
	for {
		select {
		case wa.ActiveC <- struct{}{}:
		default:
			return
		}
	}
}

// UpdateNow forces an update of the database as soon as possible.
func (wa *ChatAPIDB) UpdateNow() {
	select {
	case wa.UpdateC <- struct{}{}:
	default:
	}
}

// CopyNewMessagesLoop runs CopyNewMessages every ReadMessageInterval seconds.
func (wa *ChatAPIDB) CopyNewMessagesLoop(ctx context.Context) {
	ticker := time.NewTicker(wa.ReadMessageInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-wa.UpdateC:
			// UpdateNow was called.
		case <-ticker.C:
			// Wait for activity on the site before running an automatic update.
			select {
			case <-ctx.Done():
				return
			case <-wa.ActiveC:
			case <-wa.UpdateC:
			}
		}

		err := wa.CopyNewMessages(ctx)
		if err != nil {
			log.Printf("CopyNewMessages: %v", err)
		}
	}
}

// CopyNewMessages copies new messages from Chat-API to the database.
func (wa *ChatAPIDB) CopyNewMessages(ctx context.Context) error {
	// Request new messages from Chat-API.
	messages, number, err := wa.ChatAPI.GetMessages(ctx, GetMessagesOptions{LastMessageNumber: wa.LastMessageNumber})
	if err != nil {
		return err
	}
	if number > wa.LastMessageNumber {
		wa.LastMessageNumber = number
	}

	// Send new messages to the database.
	for _, message := range messages {
		err = wa.DB.AddMessage(ctx, "REPLACE", message)
		if err != nil {
			return err
		}
		if message.Number > wa.LastMessageNumber {
			wa.LastMessageNumber = message.Number
		}
	}

	// All new messages copied.
	return nil
}

// CopyChatMessages copies messages with chatID from Chat-API to the database.
func (wa *ChatAPIDB) CopyChatMessages(ctx context.Context, chatID string) error {
	// Request messages from one chat from Chat-API.
	messages, _, err := wa.ChatAPI.GetMessages(ctx, GetMessagesOptions{ChatID: chatID})
	if err != nil {
		return err
	}

	// Send messages to the database.
	var firstErr error
	for _, message := range messages {
		err = wa.DB.AddMessage(ctx, "REPLACE", message)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return firstErr
	}

	// All chat messages copied.
	return nil
}

func (wa *ChatAPIDB) CopyChatsLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		err := wa.CopyChats(ctx)
		if err != nil {
			log.Printf("CopyChats: %v", err)
		}
	}
}

func (wa *ChatAPIDB) CopyChats(ctx context.Context) error {
	// Request chats from Chat-API.
	chats, err := wa.ChatAPI.GetChats(ctx)
	if err != nil {
		return err
	}

	// Send chats to the database.
	var firstErr error
	for _, chat := range chats {
		err = wa.DB.AddChat(ctx, chat)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

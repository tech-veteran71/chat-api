package main

import (
	"encoding/json"
	"fmt"
)

// Chat is a dialog in Chat-API.
type Chat struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	JSON  BJSON  `json:"-"`
}

type ChatError struct {
	Chat *Chat
	Text string
}

func (err *ChatError) Error() string {
	return fmt.Sprintf("%v: %s", err.Text, err.Chat.JSON)
}

func NewChatFromBJSON(b BJSON) (*Chat, error) {
	chat := &Chat{JSON: b}

	// Decode JSON string.
	err := json.Unmarshal(b, &chat)
	if err != nil {
		return nil, err
	}

	// Check for ID.
	if chat.ID == "" {
		return nil, &ChatError{chat, "chat ID is missing"}
	}

	// Return new chat.
	return chat, nil
}

func NewChatsFromBJSON(bs []BJSON) ([]*Chat, error) {
	var firstErr error
	chats := make([]*Chat, 0, len(bs))

	for _, b := range bs {
		chat, err := NewChatFromBJSON(b)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		chats = append(chats, chat)
	}

	return chats, firstErr
}

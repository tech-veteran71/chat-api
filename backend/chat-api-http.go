// Links a Chat-API+Database to the web.

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ChatAPIHTTP implements HTTP handlers for Chat-API+Database.
type ChatAPIHTTP struct {
	*ChatAPIDB
}

func NewChatAPIHTTP(db *ChatAPIDB) *ChatAPIHTTP {
	return &ChatAPIHTTP{db}
}

type WebhookRequest struct {
	ACKs     []*WebhookACKRequest `json:"ack"`
	Messages []BJSON              `json:"messages"`
}

type WebhookACKRequest struct {
	ChatID      string `json:"chatId"`
	MessageID   string `json:"id"`
	Status      string `json:"status"`
	QueueNumber int    `json:"queueNumber"`
}

// Webhook handles incoming messages from Chat-API.
func (wa *ChatAPIHTTP) Webhook(w http.ResponseWriter, r *http.Request) {
	// Read request body.
	mb := http.MaxBytesReader(w, r.Body, 100_000_000)
	b, err := io.ReadAll(mb)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Log request.
	log.Printf("Webhook: %s", b)

	// Decode request body.
	var req WebhookRequest
	err = json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Convert BJSON's into Message's.
	messages, err := NewMessagesFromBJSON(req.Messages)
	if err != nil {
		log.Printf("NewMessagesFromBJSON: %v", err)
		// Do not return!
		// Use Message's that were successfully converted.
	}

	// Update database.
	for _, message := range messages {
		err = wa.DB.AddMessage(r.Context(), "INSERT", message)
		if err != nil && !strings.Contains(err.Error(), "UNIQUE constraint failed") {
			log.Printf("Database.AddMessage: %v", err)
			// Do not return!
			// Keep working on other messages.
		}
	}

	// Apply ACK updates.
	for _, ack := range req.ACKs {
		err = wa.DB.SetMessageAck(r.Context(), ack.ChatID, ack.MessageID, ack.Status)
		if err != nil {
			if err == sql.ErrNoRows {
				// Message is not here yet, update now.
				log.Printf("Database.SetMessageAck: Message ID not found, updating now")
				wa.UpdateNow()
			} else {
				log.Printf("Database.SetMessageAck: %v", err)
			}
			// Do not return!
			// Keep working on other updates.
		}
	}
}

// GetUserChatInfo fetches information about a user's chats.
func (wa *ChatAPIHTTP) GetUserChatInfo(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context.
	user, ok := ContextUser(r.Context())
	if !ok {
		http.Error(w, "Unknown user", http.StatusInternalServerError)
		return
	}

	// Get chat info from database.
	infos, err := wa.DB.GetUserChatInfo(r.Context(), user.ID)
	if err != nil {
		log.Printf("Database.GetUserChatInfo: %v", err)
		http.Error(w, "Cannot access database", http.StatusInternalServerError)
		return
	}

	// Send info to user.
	json.NewEncoder(w).Encode(map[string]interface{}{"infos": infos})
}

type SetUserChatAsReadRequest struct {
	ChatID    string `json:"chatID"`
	Timestamp int64  `json:"timestamp"`
}

// SetUserChatAsRead marks a chat as read.
func (wa *ChatAPIHTTP) SetUserChatAsRead(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context.
	user, ok := ContextUser(r.Context())
	if !ok {
		http.Error(w, "Unknown user", http.StatusInternalServerError)
		return
	}

	// Read request body.
	var req SetUserChatAsReadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update database.
	err = wa.DB.SetUserChatAsRead(r.Context(), user.ID, req.ChatID, req.Timestamp)
	if err != nil {
		log.Printf("Database.SetUserChatAsRead: %v", err)
		http.Error(w, "Cannot access database", http.StatusInternalServerError)
		return
	}
}

// Messages fetches messages from the database.
// New messages are not immediately fetched from Chat-API.
//
// Query parameters:
// id: only messages after this one will be fetched.
// wait: if not empty, wait for new messages to arrive.
func (wa *ChatAPIHTTP) Messages(w http.ResponseWriter, r *http.Request) {
	var err error
	uq := r.URL.Query()

	// Tell copy goroutine to keep updating.
	wa.Active()

	// Get row ID of last sent message.
	var qID int64
	if uq.Has("id") {
		qID, err = strconv.ParseInt(uq.Get("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	}

	// Create wait context.
	wait := uq.Has("wait")
	var waitCtx context.Context
	if wait {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
	}
retry:

	// Get messages from database.
	var rows []MessageRow
	if qID == 0 {
		rows, err = wa.DB.GetRecentMessages(r.Context())
		if err != nil {
			log.Printf("Database.GetRecentMessages: %v", err)
		}
	} else {
		rows, err = wa.DB.GetMessagesAfterID(r.Context(), qID)
		if err != nil {
			log.Printf("Database.GetMessagesAfterID(%v): %v", qID, err)
		}
	}
	if err != nil {
		http.Error(w, "Cannot access database", http.StatusInternalServerError)
		return
	}

	// If there are no messages, wait for new messages.
	if wait && len(rows) <= 0 {
		if wa.DB.MessageW.Wait(waitCtx) == nil {
			goto retry
		}
	}

	// Convert MessageRow's into Message's.
	messages, err := NewMessagesFromRow(rows)
	if err != nil {
		log.Printf("NewMessagesFromRow: %v", err)
		// Do not return!
		// Send messages that were successfully converted.
	}

	// Send messages to user.
	json.NewEncoder(w).Encode(map[string]interface{}{"messages": messages})
}

// MessagesByChatID fetches messages from the database for a single chat.
//
// Query parameters:
// chat_id: only messages in this chat will be fetched.
// update: synchronize database messages with Chat-API.
func (wa *ChatAPIHTTP) MessagesByChatID(w http.ResponseWriter, r *http.Request) {
	uq := r.URL.Query()

	// Get chat ID from URL.
	chatID := uq.Get("chat_id")
	if chatID == "" {
		http.Error(w, "Missing chat_id", http.StatusBadRequest)
		return
	}

	// Update database.
	if uq.Has("update") {
		err := wa.CopyChatMessages(r.Context(), chatID)
		if err != nil {
			http.Error(w, "Cannot copy messages from Chat-API", http.StatusInternalServerError)
			return
		}
	}

	// Fetch messages from database.
	rows, err := wa.DB.GetChatMessages(r.Context(), chatID)
	if err != nil {
		http.Error(w, "Cannot access database", http.StatusInternalServerError)
		return
	}

	// Convert MessageRow's into Message's.
	messages, err := NewMessagesFromRow(rows)
	if err != nil {
		log.Printf("NewMessagesFromRow: %v", err)
		// Do not return!
		// Send messages that were successfully converted.
	}

	// Send messages to user.
	json.NewEncoder(w).Encode(map[string]interface{}{"messages": messages})
}

// Chats fetches chats from the database.
// New chats are not immediately fetched from Chat-API.
//
// Query parameters:
// id: only chats after this one will be fetched.
// wait: if not empty, wait for new chats to arrive.
func (wa *ChatAPIHTTP) Chats(w http.ResponseWriter, r *http.Request) {
	var err error
	uq := r.URL.Query()

	// Get row ID of last sent chat.
	var qID int64
	if uq.Has("id") {
		qID, err = strconv.ParseInt(uq.Get("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	}

	// Create wait context.
	wait := uq.Has("wait")
	var waitCtx context.Context
	if wait {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
	}
retry:

	// Get chats from database.
	rows, err := wa.DB.GetChatsAfterID(r.Context(), qID)
	if err != nil {
		log.Printf("Database.GetChatsAfterID(%v): %v", qID, err)
		http.Error(w, "Cannot access database", http.StatusInternalServerError)
		return
	}

	// If there are no chats, wait for new chats.
	if wait && len(rows) <= 0 {
		if wa.DB.ChatW.Wait(waitCtx) == nil {
			goto retry
		}
	}

	// Convert ChatRow's into Chat's.
	chats, err := NewChatsFromRow(rows)
	if err != nil {
		log.Printf("NewChatsFromRow: %v", err)
		// Do not return!
		// Send chats that were successfully converted.
	}

	// Send chats to user.
	json.NewEncoder(w).Encode(map[string]interface{}{"chats": chats})
}

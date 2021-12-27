package main

import (
	"encoding/json"
	"fmt"
)

// Message is a chat message.
type Message struct {
	ChatID    string `json:"chatId"`
	MessageID string `json:"id"`
	Timestamp int64  `json:"time"`
	Number    int64  `json:"messageNumber"`
	JSON      BJSON  `json:"-"`
}

type MessageError struct {
	Message *Message
	Text    string
}

func (err *MessageError) Error() string {
	return fmt.Sprintf("%v: %s", err.Text, err.Message.JSON)
}

func NewMessageFromBJSON(b BJSON) (*Message, error) {
	message := &Message{JSON: b}

	// Decode JSON string.
	err := json.Unmarshal(b, &message)
	if err != nil {
		return nil, err
	}

	// Validate message.
	if message.MessageID == "" {
		return message, &MessageError{message, "message ID is missing"}
	}
	if message.ChatID == "" {
		return message, &MessageError{message, "chat ID is missing"}
	}

	// Return new message.
	return message, nil
}

func NewMessagesFromBJSON(bs []BJSON) ([]*Message, error) {
	var firstErr error
	messages := make([]*Message, 0, len(bs))

	for _, b := range bs {
		message, err := NewMessageFromBJSON(b)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		messages = append(messages, message)
	}

	return messages, firstErr
}

// NewMessageFromRow creates a new Message from a MessageRow.
func NewMessageFromRow(row MessageRow) (*Message, error) {
	message, err := NewMessageFromBJSON(row.JSON)
	if err != nil {
		return message, err
	}

	err = message.JSON.Update(func(j map[string]interface{}) error {
		j["__rowID"] = row.ID
		return nil
	})
	return message, err
}

// NewMessagesFromRow creates a slice of Message's from a slice of MessageRow's.
func NewMessagesFromRow(rows []MessageRow) ([]*Message, error) {
	var firstErr error
	messages := make([]*Message, 0, len(rows))

	for _, row := range rows {
		message, err := NewMessageFromRow(row)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		messages = append(messages, message)
	}

	return messages, firstErr
}

func (message *Message) MarshalJSON() ([]byte, error) {
	return message.JSON.MarshalJSON()
}

package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	nextMessageID int64 = 1
)

// Message is a chat message.
type Message struct {
	Number int64 // = messageNumber from the Chat-API
	Raw    map[string]interface{}
}

// NewMessageFromMap creates a new Message from a map built by the JSON package.
func NewMessageFromMap(raw map[string]interface{}) (*Message, error) {
	message := &Message{Raw: raw}

	// Copy message number.
	number, err := message.GetMessageNumber()
	if err != nil {
		return nil, err
	}
	message.Number = number

	// Return new message.
	return message, nil
}

// NewMessagesFromMap creates a slice of Message's from a slice of JSON maps.
func NewMessagesFromMap(raws []map[string]interface{}) ([]*Message, error) {
	var firstErr error
	messages := make([]*Message, 0, len(raws))

	for _, raw := range raws {
		message, err := NewMessageFromMap(raw)
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

// NewMessageFromJSON creates a new Message from a JSON string.
func NewMessageFromJSON(js string) (*Message, error) {
	// Decode JSON string.
	var raw map[string]interface{}
	err := json.Unmarshal([]byte(js), &raw)
	if err != nil {
		return nil, err
	}

	// Create message.
	return NewMessageFromMap(raw)
}

// NewMessagesFromJSON creates a slice of Message's from a slice of JSON strings.
func NewMessagesFromJSON(jss []string) ([]*Message, error) {
	var firstErr error
	messages := make([]*Message, 0, len(jss))

	for _, js := range jss {
		message, err := NewMessageFromJSON(js)
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
	message, err := NewMessageFromJSON(row.JSON)
	message.Raw["__rowID"] = row.ID
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

// NewMessage creates a new Message with a unique id.
func NewMessage(raw map[string]interface{}) (*Message, error) {
	// Unique message ID.
	atomic.AddInt64(&nextMessageID, 1)
	now := time.Now()

	// Set standard fields if unset.
	if _, ok := raw["number"]; !ok {
		raw["number"] = 0
	}
	if _, ok := raw["time"]; !ok {
		raw["time"] = float64(now.Unix())
	}
	if _, ok := raw["id"]; !ok {
		raw["id"] = fmt.Sprintf("local-%v-%v", now.UnixNano(), nextMessageID)
	}

	// Create the message.
	return NewMessageFromMap(raw)
}

func (message *Message) GetMessageNumber() (int64, error) {
	// Look for messageNumber key.
	v, ok := message.Raw["messageNumber"]
	if !ok {
		return 0, fmt.Errorf("missing key messageNumber")
	}

	// JSON decodes numbers as float64.
	number, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("messageNumber has wrong type")
	}

	// Convert to int64.
	return int64(number), nil
}

func (message *Message) GetChatID() string {
	// Look for chatID key.
	v, ok := message.Raw["chatId"]
	if !ok {
		return ""
	}

	// Must be a string.
	id, ok := v.(string)
	if !ok {
		return ""
	}

	// Success.
	return id
}

func (message *Message) MarshalJSON() ([]byte, error) {
	// The Raw map includes all the message data.
	return json.Marshal(message.Raw)
}

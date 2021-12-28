// Talks to the Chat-API service.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ChatAPI struct {
	URL   *url.URL
	Token string
}

type GetStatusResponse struct {
	Status string `json:"accountStatus"`
}

// GetStatus checks whether the account is active.
func (wa *ChatAPI) GetStatus(ctx context.Context) error {
	log.Printf("ChatAPI.GetStatus()")

	// Prepare URL.
	u := *wa.URL
	u.Path += "/status"
	q := u.Query()
	q.Add("token", wa.Token)
	u.RawQuery = q.Encode()

	// Create request.
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "https://github.com/andre-luiz-dos-santos/chat-api")

	// Send request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Read response body.
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Log response.
	log.Printf("Chat-API /status response: %s", b)

	// Decode response body.
	var j GetStatusResponse
	err = json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	// Check response.
	if j.Status != "authenticated" {
		return fmt.Errorf("not authenticated to Chat-API: %v", j.Status)
	}

	// Authenticated.
	return nil
}

type SetWebhookResponse struct {
	Set bool `json:"set"`
}

// SetWebhook sets the webhook URL.
func (wa *ChatAPI) SetWebhook(ctx context.Context, url string) error {
	log.Printf("ChatAPI.SetWebhook(%v)", url)

	// Prepare URL.
	u := *wa.URL
	u.Path += "/webhook"
	q := u.Query()
	q.Add("token", wa.Token)
	u.RawQuery = q.Encode()

	// Prepare request body.
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(map[string]interface{}{"webhookUrl": url})
	if err != nil {
		return err
	}

	// Create request.
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "https://github.com/andre-luiz-dos-santos/chat-api")
	req.Header.Set("Content-Type", "application/json")

	// Send request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Read response body.
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Log response.
	log.Printf("Chat-API /webhook response: %s", b)

	// Decode response body.
	var j SetWebhookResponse
	err = json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	// Check response.
	if !j.Set {
		return fmt.Errorf("failed to set webhook")
	}

	// Webhook updated.
	return nil
}

type GetMessagesOptions struct {
	LastMessageNumber int64
	Limit             int
	ChatID            string
	MinTime           int64
	MaxTime           int64
}

type GetMessagesResponse struct {
	LastMessageNumber *int64   `json:"lastMessageNumber"`
	Messages          *[]BJSON `json:"messages"`
	Error             string   `json:"error"`
}

// GetMessages fetches messages from Chat-API.
// Messages without messageNumber are ignored!
func (wa *ChatAPI) GetMessages(ctx context.Context, options GetMessagesOptions) ([]*Message, int64, error) {
	log.Printf("ChatAPI.GetMessages(%+v)", options)

	// Prepare URL.
	u := *wa.URL
	u.Path += "/messages"
	q := u.Query()
	q.Add("token", wa.Token)
	q.Add("limit", strconv.Itoa(options.Limit))
	if options.ChatID != "" {
		q.Add("chatId", options.ChatID)
	}
	if options.MinTime > 0 {
		q.Add("min_time", strconv.FormatInt(options.MinTime, 10))
	}
	if options.MaxTime > 0 {
		q.Add("max_time", strconv.FormatInt(options.MaxTime, 10))
	}
	if options.LastMessageNumber > 0 {
		q.Add("lastMessageNumber", strconv.FormatInt(options.LastMessageNumber, 10))
	}
	u.RawQuery = q.Encode()

	// Create request.
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "https://github.com/andre-luiz-dos-santos/chat-api")

	// Send request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	// Read response body.
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	// Log response.
	log.Printf("Chat-API /messages response: %s", b)

	// Decode response body.
	var j GetMessagesResponse
	err = json.Unmarshal(b, &j)
	if err != nil {
		return nil, 0, err
	}

	// Check response.
	if j.Error != "" {
		return nil, 0, fmt.Errorf("Chat-API /messages error: %v", j.Error)
	}
	if j.Messages == nil {
		return nil, 0, fmt.Errorf("Chat-API /messages is missing key messages")
	}
	var lastMessageNumber int64
	if j.LastMessageNumber != nil {
		lastMessageNumber = *j.LastMessageNumber
	}

	// Convert BJSON's into Message's.
	messages, err := NewMessagesFromBJSON(*j.Messages)
	if err != nil {
		log.Printf("WARNING: Ignoring message from Chat-API!")
		log.Printf("NewMessagesFromBJSON: %v", err)
	}

	// Message's fetched from Chat-API.
	return messages, lastMessageNumber, nil
}

type GetChatsResponse struct {
	Chats *[]BJSON `json:"dialogs"`
}

func (wa *ChatAPI) GetChats(ctx context.Context) ([]*Chat, error) {
	log.Printf("ChatAPI.GetChats()")

	// Prepare URL.
	u := *wa.URL
	u.Path += "/dialogs"
	q := u.Query()
	q.Add("token", wa.Token)
	u.RawQuery = q.Encode()

	// Create request.
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "https://github.com/andre-luiz-dos-santos/chat-api")

	// Send request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response body.
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Log response.
	log.Printf("Chat-API /dialogs response: %s", b)

	// Decode response body.
	var j GetChatsResponse
	err = json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	// Check response.
	if j.Chats == nil {
		return nil, fmt.Errorf("Chat-API /dialogs is missing key dialogs")
	}

	// Create Chat instances.
	chats, err := NewChatsFromBJSON(*j.Chats)
	if err != nil {
		log.Printf("WARNING: Ignoring chat from Chat-API!")
		log.Printf("NewChatsFromBJSON: %v", err)
	}

	// Chat's fetched from Chat-API.
	return chats, nil
}

// ackToNum converts an ack string to a comparable number.
func ackToNum(ack string) int {
	switch ack {
	case "viewed":
		return 10
	case "read":
		return 9
	case "delivered":
		return 7
	case "sent":
		return 5
	default:
		return 0
	}
}

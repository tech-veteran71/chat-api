// Manage database connection, read and write to the database.

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Driver string
	DSN    string

	*sql.DB
	// TODO: Check why two goroutines writing to the database at the same time fails.
	sync.Mutex

	ChatW    Wait
	MessageW Wait
}

func NewDatabase(driver, dsn string) *Database {
	return &Database{
		Driver: driver,
		DSN:    dsn,
	}
}

// Open opens the database.
func (db *Database) Open() error {
	h, err := sql.Open(db.Driver, db.DSN)
	if err != nil {
		return err
	}
	db.DB = h

	return nil
}

// Create creates tables, indexes, etc.
func (db *Database) Create() error {
	_, err := db.Exec(SQLITE_INIT)
	return err
}

// CheckPassword returns User if username and password exist in the users table.
func (db *Database) CheckPassword(ctx context.Context, username, password string) (*User, error) {
	var user User

	err := db.QueryRowContext(ctx, `SELECT id, name, label, token FROM users WHERE name = ? AND password = ? LIMIT 1`, username, password).
		Scan(&user.ID, &user.Name, &user.Label, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidUser
		}
		return nil, err
	}

	return &user, nil
}

// CheckToken returns User if token exists in the users table.
func (db *Database) CheckToken(ctx context.Context, token string) (*User, error) {
	var user User

	err := db.QueryRowContext(ctx, `SELECT id, name, label FROM users WHERE token = ? LIMIT 1`, token).
		Scan(&user.ID, &user.Name, &user.Label)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidUser
		}
		return nil, err
	}

	user.Token = token
	return &user, nil
}

// ResetToken resets the token for username.
func (db *Database) ResetToken(ctx context.Context, username, token string) error {
	res, err := db.ExecContext(ctx, `UPDATE users SET token = ? WHERE name = ?`, token, username)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrInvalidUser
	}

	return nil
}

type UserChatInfo struct {
	ChatID     string `json:"chatID"`
	ReadBefore int64  `json:"readBefore"`
}

func (db *Database) GetUserChatInfo(ctx context.Context, userID int64) ([]*UserChatInfo, error) {
	infos := make([]*UserChatInfo, 0, 128)

	rows, err := db.QueryContext(ctx, `SELECT chat_id, read_before FROM user_chat_info WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info UserChatInfo
		err = rows.Scan(&info.ChatID, &info.ReadBefore)
		if err != nil {
			return nil, err
		}
		infos = append(infos, &info)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (db *Database) SetUserChatAsRead(ctx context.Context, userID int64, chatID string, timestamp int64) error {
	db.Lock()
	defer db.Unlock()

	res, err := db.ExecContext(ctx, `UPDATE user_chat_info SET read_before = ? WHERE user_id = ? AND chat_id = ?`, timestamp, userID, chatID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected >= 1 {
		// UPDATE succeeded.
		return nil
	}

	_, err = db.ExecContext(ctx, `INSERT INTO user_chat_info (user_id, chat_id, read_before) VALUES (?, ?, ?)`, userID, chatID, timestamp)
	if err != nil {
		return err
	}

	// INSERT succeeded.
	return nil
}

// AddMessage adds Message to the database.
func (db *Database) AddMessage(ctx context.Context, cmd string, message *Message) error {
	db.Lock()
	defer db.Unlock()

	// Return if message is already in the database.
	var id int64
	err := db.QueryRowContext(ctx, `SELECT id FROM messages WHERE message_id = ? AND json = ?`, message.ID, message.JSON).Scan(&id)
	if err == nil || err != sql.ErrNoRows {
		return err
	}

	// INSERT or REPLACE new message (needs new id).
	_, err = db.ExecContext(ctx,
		cmd+` INTO messages (time, message_number, message_id, chat_id, json) VALUES (?, ?, ?, ?, ?)`,
		message.Timestamp, message.Number, message.ID, message.ChatID, message.JSON)
	if err != nil {
		return err
	}

	db.MessageW.Notify()
	return nil
}

func (db *Database) AddChat(ctx context.Context, chat *Chat) error {
	db.Lock()
	defer db.Unlock()

	// Return if chat is already in the database.
	var id int64
	err := db.QueryRowContext(ctx, `SELECT id FROM chats WHERE chat_id = ? AND json = ? LIMIT 1`, chat.ID, chat.JSON).Scan(&id)
	if err == nil || err != sql.ErrNoRows {
		return err
	}

	// INSERT or REPLACE new chat (needs new id).
	_, err = db.ExecContext(ctx,
		`REPLACE INTO chats (chat_id, json) VALUES (?, ?)`,
		chat.ID, chat.JSON)
	if err != nil {
		return err
	}

	db.ChatW.Notify()
	return nil
}

// SetMessageAck sets the ack field of a Message.
// The row ID must be incremented for the browser to detect the change.
func (db *Database) SetMessageAck(ctx context.Context, chatID, messageID, ack string) error {
	db.Lock()
	defer db.Unlock()

	// Read current row.
	var time, messageNumber int64
	var js string // JSON string
	err := db.QueryRowContext(ctx, `SELECT time, message_number, json FROM messages WHERE chat_id = ? AND message_id = ? LIMIT 1`, chatID, messageID).
		Scan(&time, &messageNumber, &js)
	if err != nil {
		return err
	}

	// Decode JSON.
	var raw map[string]interface{}
	err = json.Unmarshal([]byte(js), &raw)
	if err != nil {
		return err
	}

	// Update ack if it's newer.
	if dbValue, ok := raw["ack"]; ok {
		if dbAck, ok := dbValue.(string); ok {
			if ackToNum(dbAck) > ackToNum(ack) {
				log.Printf("Ignoring ack update: %v > %v", dbAck, ack)
				return nil
			}
		}
	}
	raw["ack"] = ack

	// Encode message.
	b, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	js = string(b)

	// REPLACE row.
	_, err = db.ExecContext(ctx, `REPLACE INTO messages (time, message_number, message_id, chat_id, json) VALUES (?, ?, ?, ?, ?)`, time, messageNumber, messageID, chatID, js)
	if err != nil {
		return err
	}

	db.MessageW.Notify()
	return nil
}

// GetLastMessageNumber returns the largest message number in the database.
func (db *Database) GetLastMessageNumber(ctx context.Context) (int64, error) {
	var number int64

	err := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(message_number), 0) FROM messages`).
		Scan(&number)
	if err != nil {
		return 0, err
	}

	return number, nil
}

type MessageRow struct {
	ID   int64
	JSON []byte
}

// GetRecentMessages returns the last 1000 messages ordered by time.
func (db *Database) GetRecentMessages(ctx context.Context) ([]MessageRow, error) {
	var messages []MessageRow

	sql_chats := `SELECT *, row_number() OVER (PARTITION BY chat_id ORDER BY time DESC) AS row FROM messages`
	rows, err := db.QueryContext(ctx, `SELECT id, json FROM (`+sql_chats+`) AS messages WHERE row <= 20 ORDER BY time DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var json []byte
		err = rows.Scan(&id, &json)
		if err != nil {
			return nil, err
		}
		messages = append(messages, MessageRow{id, json})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessagesAfterID returns messages after id.
func (db *Database) GetMessagesAfterID(ctx context.Context, id int64) ([]MessageRow, error) {
	var messages []MessageRow

	rows, err := db.QueryContext(ctx, `SELECT id, json FROM messages WHERE id > ? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var json []byte
		err = rows.Scan(&id, &json)
		if err != nil {
			return nil, err
		}
		messages = append(messages, MessageRow{id, json})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetChatMessages returns messages in chatID.
func (db *Database) GetChatMessages(ctx context.Context, chatID string) ([]MessageRow, error) {
	var messages []MessageRow

	rows, err := db.QueryContext(ctx, `SELECT id, json FROM messages WHERE chat_id = ?`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var json []byte
		err = rows.Scan(&id, &json)
		if err != nil {
			return nil, err
		}
		messages = append(messages, MessageRow{id, json})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

type ChatRow struct {
	ID   int64
	JSON []byte
}

func (db *Database) GetChatsAfterID(ctx context.Context, id int64) ([]ChatRow, error) {
	var chats []ChatRow

	rows, err := db.QueryContext(ctx, `SELECT id, json FROM chats WHERE id > ? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var json []byte
		err = rows.Scan(&id, &json)
		if err != nil {
			return nil, err
		}
		chats = append(chats, ChatRow{id, json})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return chats, nil
}

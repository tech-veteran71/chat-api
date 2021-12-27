package main

var SQLITE_INIT = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	name TEXT,
	label TEXT,
	password TEXT,
	token TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON users (name);

-- TODO: User management.
INSERT OR IGNORE INTO users (name, label, password, token) VALUES ("admin", "Admin", "admin", "");

CREATE TABLE IF NOT EXISTS user_chat_info (
	id INTEGER PRIMARY KEY,
	user_id INTEGER, -- references users.id
	chat_id TEXT, -- references messages.chat_id
	read_before INTEGER -- time of last read message
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_chat_info_ids ON user_chat_info (user_id, chat_id);

CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY,
	time INTEGER,
	message_number INTEGER,
	message_id TEXT,
	chat_id TEXT,
	json TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_messages_id ON messages (message_id);

CREATE TABLE IF NOT EXISTS chats (
	id INTEGER PRIMARY KEY,
	chat_id TEXT,
	json TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_chats_id ON chats (chat_id);
`

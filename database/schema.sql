-- =========================================================
-- CHAT PLATFORM DATABASE SCHEMA
-- =========================================================

PRAGMA foreign_keys = ON;

-- =========================================================
-- USERS
-- =========================================================

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL,

    username TEXT NOT NULL UNIQUE,
    email    TEXT NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,

    avatar TEXT DEFAULT '',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- =========================================================
-- SESSIONS
-- =========================================================

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    user_id INTEGER NOT NULL,

    token TEXT NOT NULL UNIQUE,

    expires_at DATETIME NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================================================
-- MESSAGES
-- =========================================================

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    -- User who sent the message
    sender_id INTEGER NOT NULL,

    -- NULL means public chat
    receiver_id INTEGER,

    -- Text message or media caption
    message TEXT,

    -- text | image | video | document
    message_type TEXT NOT NULL DEFAULT 'text',

    -- Optional uploaded file
    upload_id INTEGER,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (sender_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    FOREIGN KEY (receiver_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================================================
-- FILE UPLOADS
-- =========================================================

CREATE TABLE IF NOT EXISTS uploads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    file_name TEXT NOT NULL,

    file_path TEXT NOT NULL,

    file_size INTEGER NOT NULL,

    mime_type TEXT NOT NULL,

    uploaded_by INTEGER NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (uploaded_by)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================================================
-- INDEXES
-- =========================================================

CREATE INDEX IF NOT EXISTS idx_users_email
ON users(email);

CREATE INDEX IF NOT EXISTS idx_users_username
ON users(username);

CREATE INDEX IF NOT EXISTS idx_sessions_token
ON sessions(token);

CREATE INDEX IF NOT EXISTS idx_messages_sender
ON messages(sender_id);

CREATE INDEX IF NOT EXISTS idx_messages_receiver
ON messages(receiver_id);

CREATE INDEX IF NOT EXISTS idx_messages_created_at
ON messages(created_at);

CREATE INDEX IF NOT EXISTS idx_uploads_user
ON uploads(uploaded_by);
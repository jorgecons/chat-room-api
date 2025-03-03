CREATE TABLE IF NOT EXISTS users
(
    username
    TEXT
    PRIMARY
    KEY,
    password_hash
    TEXT
    NOT
    NULL
);

CREATE TABLE IF NOT EXISTS messages
(
    id
    SERIAL
    PRIMARY
    KEY,
    room
    TEXT
    NOT
    NULL,
    username
    TEXT
    NOT
    NULL,
    text
    TEXT
    NOT
    NULL,
    date
    TIMESTAMP
    DEFAULT
    NOW
(
),
    FOREIGN KEY
(
    username
) REFERENCES users
(
    username
) ON DELETE CASCADE
    );

-- Indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_messages_room ON messages(room);
CREATE INDEX IF NOT EXISTS idx_messages_date ON messages(date DESC);

-- Ensure unique usernames
ALTER TABLE users
    ADD CONSTRAINT unique_username UNIQUE (username);

-- Ensure the 'bot' user exists. password is "bot_password"
INSERT INTO users (username, password_hash)
VALUES ('bot', '$2b$12$ctWZ.S8ognRCBpQ.Cr7qbum27A4z3ShC/5rBzOgEk.hCem62vfkGy') ON CONFLICT (username) DO NOTHING;

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    username TEXT,
    text TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,

    CONSTRAINT unique_message
        UNIQUE(chat_id, message_id)
);

CREATE TABLE users (
    id BIGINT PRIMARY KEY, -- Telegram user ID
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chats (
    id BIGINT PRIMARY KEY, -- Telegram chat ID
    title VARCHAR(255),
    type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ONLY admins stored here
CREATE TABLE chat_admins (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    role TEXT NOT NULL CHECK (role IN ('creator', 'administrator')),

    PRIMARY KEY (chat_id, user_id),

    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE bot_configs (
    chat_id BIGINT PRIMARY KEY,

    summary_enabled BOOLEAN DEFAULT FALSE,
    duplicate_dm_enabled BOOLEAN DEFAULT FALSE,

    updated_by BIGINT,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE command_configs (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,

    command_name VARCHAR(50) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,

    UNIQUE (chat_id, command_name),

    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);

-- Maintain updated_at automatically (Postgres replacement for ON UPDATE)
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_bot_configs_updated_at ON bot_configs;
CREATE TRIGGER trg_bot_configs_updated_at
BEFORE UPDATE ON bot_configs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
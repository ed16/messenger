-- +goose Up
CREATE TYPE message_status AS ENUM ('sent', 'received', 'read');
CREATE TABLE messages (
    message_id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    recipient_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    content VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status message_status,
    media_id BIGINT,
    CONSTRAINT fk_media FOREIGN KEY (media_id) REFERENCES media(media_id) ON DELETE RESTRICT
);
CREATE INDEX idx_sender_id ON messages (sender_id);

-- +goose Down
DROP TABLE IF EXISTS messages CASCADE;
DROP INDEX IF EXISTS idx_sender_id;
DROP TYPE IF EXISTS message_status;
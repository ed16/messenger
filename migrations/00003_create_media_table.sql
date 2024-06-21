-- +goose Up
CREATE TABLE media (
    media_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    file_path VARCHAR(260) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_media_user_id ON media (user_id);

-- +goose Down
DROP TABLE IF EXISTS media CASCADE;
DROP INDEX IF EXISTS idx_user_id;
-- +goose Up
CREATE TABLE contacts (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    contact_user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, contact_user_id)
);
CREATE INDEX idx_contacts_user_id ON contacts (user_id);

-- +goose Down
DROP TABLE IF EXISTS contacts CASCADE;
DROP INDEX IF EXISTS idx_user_id;
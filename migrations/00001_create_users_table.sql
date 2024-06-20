-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_status AS ENUM ('active', 'deleted', 'blocked');
CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    status user_status,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    password_hash CHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS user_status;
-- +goose StatementEnd
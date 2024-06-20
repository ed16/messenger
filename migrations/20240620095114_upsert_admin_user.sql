-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, password_hash, status) 
VALUES ('admin', '\$2a\$10\$p7X62PHGUAGFnhdBDLFjs.ufDZY.59FbWlrBi1PxG4OKlHEb.lTVO', 'active')
ON CONFLICT (username) 
DO 
UPDATE SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE username = 'admin';
-- +goose StatementEnd

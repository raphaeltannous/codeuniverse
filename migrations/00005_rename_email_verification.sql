-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS users
RENAME COLUMN email_verified TO is_verified;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS users
RENAME COLUMN is_verified TO email_verified;
-- +goose StatementEnd

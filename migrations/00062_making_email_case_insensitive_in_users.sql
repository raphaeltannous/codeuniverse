-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN email TYPE CITEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN email TYPE VARCHAR(255);
-- +goose StatementEnd

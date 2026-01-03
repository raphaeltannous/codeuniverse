-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS CITEXT;

ALTER TABLE users
ALTER COLUMN username TYPE CITEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN username TYPE VARCHAR(50);

DROP EXTENSION IF EXISTS CITEXT;
-- +goose StatementEnd

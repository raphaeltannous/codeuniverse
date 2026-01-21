-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN premium_status TEXT NOT NULL DEFAULT 'free';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN premium_status;
-- +goose StatementEnd

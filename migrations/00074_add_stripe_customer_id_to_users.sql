-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN stripe_customer_id TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN stripe_customer_id;
-- +goose StatementEnd

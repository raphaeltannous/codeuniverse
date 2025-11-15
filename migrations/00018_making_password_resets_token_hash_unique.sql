-- +goose Up
-- +goose StatementBegin
ALTER TABLE password_resets
ADD CONSTRAINT password_resets_token_hash_key UNIQUE (token_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE password_resets
ADD CONSTRAINT password_resets_token_hash_key;
-- +goose StatementEnd

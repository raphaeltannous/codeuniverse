-- +goose Up
-- +goose StatementBegin
ALTER TABLE email_verifications
ADD CONSTRAINT email_verifications_token_hash_key UNIQUE (token_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE email_verifications
ADD CONSTRAINT email_verifications_token_hash_key;
-- +goose StatementEnd

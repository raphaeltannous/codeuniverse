-- +goose Up
-- +goose StatementBegin
ALTER TABLE mfa_codes
ADD CONSTRAINT mfa_codes_token_hash_key UNIQUE (token_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mfa_codes
ADD CONSTRAINT mfa_codes_token_hash_key;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE mfa_codes
ADD COLUMN token_hash TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mfa_codes
DROP COLUMN token_hash;
-- +goose StatementEnd

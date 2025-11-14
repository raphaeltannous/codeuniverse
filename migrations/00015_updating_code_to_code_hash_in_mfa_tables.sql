-- +goose Up
-- +goose StatementBegin
ALTER TABLE mfa_codes
RENAME COLUMN code to code_hash;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mfa_codes
RENAME COLUMN code_hash to code;
-- +goose StatementEnd

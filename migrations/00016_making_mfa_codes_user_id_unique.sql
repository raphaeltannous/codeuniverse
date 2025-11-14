-- +goose Up
-- +goose StatementBegin
ALTER TABLE mfa_codes
ADD CONSTRAINT mfa_codes_user_id_key UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mfa_codes
DROP CONSTRAINT mfa_codes_user_id_key;
-- +goose StatementEnd

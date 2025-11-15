-- +goose Up
-- +goose StatementBegin
ALTER TABLE email_verifications
ADD CONSTRAINT email_verifications_user_id_key UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE email_verifications
ADD CONSTRAINT email_verifications_user_id_key;
-- +goose StatementEnd

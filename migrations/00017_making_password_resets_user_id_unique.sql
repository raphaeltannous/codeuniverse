-- +goose Up
-- +goose StatementBegin
ALTER TABLE password_resets
ADD CONSTRAINT password_resets_user_id_key UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE password_resets
DROP CONSTRAINT password_resets_user_id_key;
-- +goose StatementEnd

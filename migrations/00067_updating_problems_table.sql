-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
ALTER COLUMN difficulty SET NOT NULL,
ALTER COLUMN is_premium SET NOT NULL,
ALTER COLUMN is_public SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
ALTER COLUMN difficulty DROP NOT NULL,
ALTER COLUMN is_premium DROP NOT NULL,
ALTER COLUMN is_public DROP NOT NULL;
-- +goose StatementEnd

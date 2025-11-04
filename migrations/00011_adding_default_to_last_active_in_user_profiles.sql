-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_profiles
    ALTER COLUMN last_active SET DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_profiles
    ALTER COLUMN last_active DROP DEFAULT;
-- +goose StatementEnd

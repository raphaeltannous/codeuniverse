-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_profiles
    ALTER COLUMN avatar_url SET DEFAULT 'default.png';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_profiles
    ALTER COLUMN avatar_url DROP DEFAULT;
-- +goose StatementEnd

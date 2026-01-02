-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_profiles
DROP COLUMN avatar_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_profiles
ADD COLUMN avatar_url TEXT DEFAULT 'default.png';
-- +goose StatementEnd

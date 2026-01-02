-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN avatar_url VARCHAR(255) DEFAULT 'default.png';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN avatar_url;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE courses
    ALTER COLUMN thumbnail_url SET DEFAULT 'default.jpg';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courses
    ALTER COLUMN thumbnail_url DROP DEFAULT;
-- +goose StatementEnd

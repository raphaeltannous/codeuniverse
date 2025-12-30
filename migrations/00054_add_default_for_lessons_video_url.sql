-- +goose Up
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN video_url SET DEFAULT 'default.mp4';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN video_url DROP DEFAULT;
-- +goose StatementEnd

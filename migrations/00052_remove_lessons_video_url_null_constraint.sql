-- +goose Up
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN video_url DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN video_url SET NOT NULL;
-- +goose StatementEnd

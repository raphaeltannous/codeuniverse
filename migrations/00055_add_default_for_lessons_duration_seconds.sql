-- +goose Up
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN duration_seconds SET DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN duration_seconds DROP DEFAULT;
-- +goose StatementEnd

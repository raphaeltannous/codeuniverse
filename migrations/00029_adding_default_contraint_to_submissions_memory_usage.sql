-- +goose Up
-- +goose StatementBegin
ALTER TABLE submissions
ALTER COLUMN memory_usage SET DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE submissions
ALTER COLUMN memory_usage DROP DEFAULT;
-- +goose StatementEnd

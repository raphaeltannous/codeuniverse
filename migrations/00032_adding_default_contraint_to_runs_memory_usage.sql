-- +goose Up
-- +goose StatementBegin
ALTER TABLE runs
ALTER COLUMN memory_usage SET DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE runs
ALTER COLUMN memory_usage DROP DEFAULT;
-- +goose StatementEnd

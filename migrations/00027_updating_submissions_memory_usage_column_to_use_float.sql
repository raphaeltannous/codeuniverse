-- +goose Up
-- +goose StatementBegin
ALTER TABLE submissions
ALTER COLUMN memory_usage TYPE FLOAT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE submissions
ALTER COLUMN memory_usage TYPE INT;
-- +goose StatementEnd

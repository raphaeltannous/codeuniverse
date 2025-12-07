-- +goose Up
-- +goose StatementBegin
ALTER TABLE submissions
ALTER COLUMN execution_time SET DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE submissions
ALTER COLUMN execution_time DROP DEFAULT;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN duration_seconds DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN duration_seconds SET NOT NULL;
-- +goose StatementEnd

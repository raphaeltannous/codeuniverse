-- +goose Up
-- +goose StatementBegin
ALTER TABLE submissions
    ADD COLUMN stdout TEXT NOT NULL DEFAULT '',
    ADD COLUMN stderr TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE submissions
    DROP COLUMN stdout,
    DROP COLUMN stderr;
-- +goose StatementEnd

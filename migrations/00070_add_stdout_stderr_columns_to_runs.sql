-- +goose Up
-- +goose StatementBegin
ALTER TABLE runs
    ADD COLUMN stdout TEXT NOT NULL DEFAULT '',
    ADD COLUMN stderr TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE runs
    DROP COLUMN stdout,
    DROP COLUMN stderr;
-- +goose StatementEnd

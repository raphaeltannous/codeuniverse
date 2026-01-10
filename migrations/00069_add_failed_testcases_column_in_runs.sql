-- +goose Up
-- +goose StatementBegin
ALTER TABLE runs
    ADD COLUMN failed_testcases JSONB NOT NULL DEFAULT '[]';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE runs
    DROP COLUMN failed_testcases;
-- +goose StatementEnd

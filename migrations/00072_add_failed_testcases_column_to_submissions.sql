-- +goose Up
-- +goose StatementBegin
ALTER TABLE submissions
    ADD COLUMN failed_testcases JSONB NOT NULL DEFAULT '[]';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE submissions
    DROP COLUMN failed_testcases;
-- +goose StatementEnd

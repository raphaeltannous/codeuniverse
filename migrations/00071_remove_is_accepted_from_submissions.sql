-- +goose Up
-- +goose StatementBegin
ALTER TABLE submissions
    DROP COLUMN is_accepted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE submissions
    ADD COLUMN is_accepted BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

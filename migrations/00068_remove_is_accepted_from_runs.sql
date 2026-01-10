-- +goose Up
-- +goose StatementBegin
ALTER TABLE runs
DROP COLUMN is_accepted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE runs
ADD COLUMN is_accepted BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

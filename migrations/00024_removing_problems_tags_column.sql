-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
DROP COLUMN tags;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
ADD COLUMN tags TEXT[];
-- +goose StatementEnd

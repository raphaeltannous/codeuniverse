-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
DROP COLUMN likes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
ADD COLUMN likes INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
DROP COLUMN dislikes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
ADD COLUMN dislikes INT DEFAULT 0;
-- +goose StatementEnd

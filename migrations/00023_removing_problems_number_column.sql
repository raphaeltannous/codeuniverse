-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
DROP COLUMN number;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
ADD COLUMN number INT NOT NULL;
-- +goose StatementEnd

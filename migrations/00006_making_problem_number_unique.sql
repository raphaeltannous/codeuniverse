-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS problems
ADD CONSTRAINT unique_number UNIQUE (number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS problems
DROP CONSTRAINT unique_number;
-- +goose StatementEnd

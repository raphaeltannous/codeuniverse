-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
ADD CONSTRAINT unique_number UNIQUE (number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
DROP CONSTRAINT unique_number;
-- +goose StatementEnd

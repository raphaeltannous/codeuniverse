-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
ALTER COLUMN hints SET DEFAULT ARRAY[]::TEXT[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
ALTER COLUMN hints DROP DEFAULT;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER problems_update_updated_at
BEFORE UPDATE ON problems
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS problems_update_updated_at ON problems;
-- +goose StatementEnd

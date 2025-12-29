-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER lessons_update_updated_at
BEFORE UPDATE ON lessons
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS lessons_update_updated_at ON lessons;
-- +goose StatementEnd

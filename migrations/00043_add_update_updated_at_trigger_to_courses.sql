-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER courses_update_updated_at
BEFORE UPDATE ON courses
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS courses_update_updated_at ON courses;
-- +goose StatementEnd

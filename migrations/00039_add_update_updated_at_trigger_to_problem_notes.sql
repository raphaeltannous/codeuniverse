-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER problem_notes_update_updated_at
BEFORE UPDATE ON problem_notes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS problem_notes_update_updated_at ON problem_notes;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER problem_reactions_update_updated_at
BEFORE UPDATE ON problem_reactions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS problem_reactions_update_updated_at ON problem_reactions;
-- +goose StatementEnd

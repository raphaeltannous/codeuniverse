-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER runs_update_updated_at
BEFORE UPDATE ON runs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS runs_update_updated_at ON runs;
-- +goose StatementEnd

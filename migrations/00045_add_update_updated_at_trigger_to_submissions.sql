-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER submissions_update_updated_at
BEFORE UPDATE ON submissions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS submissions_update_updated_at ON submissions;
-- +goose StatementEnd

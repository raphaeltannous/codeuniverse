-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER submissions_update_user_profiles_last_active_column
AFTER INSERT ON submissions
FOR EACH ROW
EXECUTE FUNCTION update_user_profiles_last_active_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS submissions_update_user_profiles_last_active_column ON submissions;
-- +goose StatementEnd

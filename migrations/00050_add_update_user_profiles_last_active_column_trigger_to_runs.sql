-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER runs_update_user_profiles_last_active_column
AFTER INSERT ON runs
FOR EACH ROW
EXECUTE FUNCTION update_user_profiles_last_active_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS runs_update_user_profiles_last_active_column ON runs;
-- +goose StatementEnd

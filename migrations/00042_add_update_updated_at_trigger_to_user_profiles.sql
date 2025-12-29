-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER user_profiles_update_updated_at
BEFORE UPDATE ON user_profiles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS user_profiles_update_updated_at ON user_profiles;
-- +goose StatementEnd

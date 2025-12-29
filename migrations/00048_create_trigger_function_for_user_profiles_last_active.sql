-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_user_profiles_last_active_column()
RETURNS trigger AS $$
BEGIN
    UPDATE user_profiles
    SET last_active = CURRENT_TIMESTAMP
    WHERE user_id = NEW.user_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS update_user_profiles_last_active_column;
-- +goose StatementEnd

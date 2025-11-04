-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_profiles
    DROP COLUMN total_submissions,
    DROP COLUMN accepted_submissions,
    DROP COLUMN problems_solved,
    DROP COLUMN easy_solved,
    DROP COLUMN medium_solved,
    DROP COLUMN hard_solved;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_profiles
    ADD COLUMN total_submissions INT DEFAULT 0,
    ADD COLUMN accepted_submissions INT DEFAULT 0,
    ADD COLUMN problems_solved INT DEFAULT 0,
    ADD COLUMN easy_solved INT DEFAULT 0,
    ADD COLUMN medium_solved INT DEFAULT 0,
    ADD COLUMN hard_solved INT DEFAULT 0;
-- +goose StatementEnd

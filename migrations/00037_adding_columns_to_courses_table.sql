-- +goose Up
-- +goose StatementBegin
ALTER TABLE courses
ADD COLUMN slug VARCHAR(50) UNIQUE,
ADD COLUMN difficulty VARCHAR(50) DEFAULT 'Beginner';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courses
DROP COLUMN slug,
DROP COLUMN difficulty;
-- +goose StatementEnd

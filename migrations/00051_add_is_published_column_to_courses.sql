-- +goose Up
-- +goose StatementBegin
ALTER TABLE courses
ADD COLUMN is_published BOOLEAN DEFAULT False;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courses
DROP COLUMN is_published;
-- +goose StatementEnd

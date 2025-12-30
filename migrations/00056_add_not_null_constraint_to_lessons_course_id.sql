-- +goose Up
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN course_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lessons
ALTER COLUMN course_id DROP NOT NULL;
-- +goose StatementEnd

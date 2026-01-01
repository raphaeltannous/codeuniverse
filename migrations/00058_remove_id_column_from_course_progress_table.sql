-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_progress
DROP COLUMN id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_progress
ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
-- +goose StatementEnd

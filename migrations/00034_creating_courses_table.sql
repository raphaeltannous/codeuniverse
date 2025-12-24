-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- course info
    title VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    thumbnail_url TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE courses;
-- +goose StatementEnd

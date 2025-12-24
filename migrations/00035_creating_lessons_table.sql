-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,

    -- course info
    title VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    video_url TEXT NOT NULL,
    duration_seconds INT NOT NULL,
    lesson_number INT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lessons;
-- +goose StatementEnd

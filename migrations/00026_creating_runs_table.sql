-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    problem_id UUID REFERENCES problems(id) ON DELETE CASCADE,

    -- run info
    language VARCHAR(50) NOT NULL,
    code TEXT NOT NULL,
    status VARCHAR(50),
    execution_time FLOAT,
    memory_usage FLOAT,
    is_accepted BOOLEAN DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE runs;
-- +goose StatementEnd

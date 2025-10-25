-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- problem information
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    number INT NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(20),
    tags TEXT[],
    hints TEXT[],
    code_snippets JSONB DEFAULT '{}',
    test_cases JSONB DEFAULT '{}',
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,

    -- visibility
    is_paid BOOLEAN DEFAULT false,
    is_public BOOLEAN DEFAULT true,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS problems;
-- +goose StatementEnd

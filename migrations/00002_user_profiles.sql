-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- profile information
    name VARCHAR(255),
    bio TEXT,
    avatar_url TEXT,
    country VARCHAR(100),
    preferred_language VARCHAR(50),

    -- coding stats
    total_submissions INT DEFAULT 0,
    accepted_submissions INT DEFAULT 0,
    problems_solved INT DEFAULT 0,
    easy_solved INT DEFAULT 0,
    medium_solved INT DEFAULT 0,
    hard_solved INT DEFAULT 0,

    -- social links
    website_url TEXT,
    github_url TEXT,
    linkedin_url TEXT,
    x_url TEXT,

    -- activity
    last_active TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd

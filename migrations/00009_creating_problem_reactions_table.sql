-- +goose Up
-- +goose StatementBegin
CREATE TYPE reaction_type
AS ENUM('like', 'dislike', 'none');

CREATE TABLE IF NOT EXISTS problem_reactions (
    reaction_id BIGSERIAL PRIMARY KEY,

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,

    reaction reaction_type NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS problem_reactions;
DROP TYPE IF EXISTS reaction_type;
-- +goose StatementEnd

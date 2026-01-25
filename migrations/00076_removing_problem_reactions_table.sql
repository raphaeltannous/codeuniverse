-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS problem_reactions_update_updated_at ON problem_reactions;

DROP TABLE problem_reactions;
DROP TYPE reaction_type;
-- +goose StatementEnd

-- +goose Down
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

CREATE TRIGGER problem_reactions_update_updated_at
BEFORE UPDATE ON problem_reactions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

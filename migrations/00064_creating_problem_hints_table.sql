-- +goose Up
-- +goose StatementBegin
CREATE TABLE problem_hints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,

    hint TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE problem_hints;
-- +goose StatementEnd

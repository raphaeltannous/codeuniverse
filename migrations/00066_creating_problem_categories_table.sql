-- +goose Up
-- +goose StatementBegin
CREATE TABLE problem_categories (
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (problem_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE problem_categories;
-- +goose StatementEnd

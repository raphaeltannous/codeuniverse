-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems
    DROP COLUMN hints,
    DROP COLUMN code_snippets,
    DROP COLUMN test_cases;

ALTER TABLE problems
    RENAME COLUMN is_paid TO is_premium;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems
    ADD COLUMN hints TEXT[] DEFAULT ARRAY[]::TEXT[],
    ADD COLUMN code_snippets JSONB DEFAULT '{}',
    ADD COLUMN test_cases JSONB DEFAULT '{}';

ALTER TABLE problems
    RENAME COLUMN is_premium TO is_paid;
-- +goose StatementEnd

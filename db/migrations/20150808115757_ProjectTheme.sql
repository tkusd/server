-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE projects ADD theme VARCHAR(8) NOT NULL DEFAULT 'ios';
UPDATE projects SET theme = 'ios' WHERE theme IS NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE projects DROP theme;

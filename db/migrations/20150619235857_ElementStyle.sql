-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE elements ADD styles JSONB NOT NULL DEFAULT '{}';
UPDATE elements SET styles = '{}' WHERE styles IS NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE elements DROP styles;
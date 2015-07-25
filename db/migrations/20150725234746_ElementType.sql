-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE elements ALTER COLUMN type TYPE VARCHAR(32);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

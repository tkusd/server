-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE elements RENAME COLUMN order_id TO index;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE elements RENAME COLUMN index TO order_id;

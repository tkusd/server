-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE elements ADD is_visible BOOLEAN NOT NULL DEFAULT TRUE;
UPDATE elements SET is_visible = TRUE WHERE is_visible IS NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE elements DROP is_visible;
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE assets ADD hash CHAR(40) NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE assets DROP hash;

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM tokens;
ALTER TABLE tokens ADD secret CHAR(64) NOT NULL UNIQUE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE tokens DROP COLUMN secret;

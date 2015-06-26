-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD language VARCHAR(35) NOT NULL DEFAULT 'en';
UPDATE users SET language = 'en' WHERE language IS NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP language;

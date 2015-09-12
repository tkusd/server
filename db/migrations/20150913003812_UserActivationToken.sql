-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP VIEW IF EXISTS users_extended;
ALTER TABLE users DROP COLUMN activation_token;
ALTER TABLE users ADD password_reset_token UUID;
ALTER TABLE users ADD password_reset_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD activation_token UUID;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP COLUMN activation_token;
ALTER TABLE users DROP COLUMN password_reset_token;
ALTER TABLE users DROP COLUMN password_reset_at;
ALTER TABLE users ADD activation_token CHAR(64);
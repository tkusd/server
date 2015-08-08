-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE assets ADD size BIGINT NOT NULL CHECK (size >= 0) DEFAULT 0;
ALTER TABLE assets ADD type VARCHAR(255) NOT NULL DEFAULT 'application/octet-stream';
ALTER TABLE assets ADD slug VARCHAR(255) NOT NULL;
ALTER TABLE assets ADD width INTEGER;
ALTER TABLE assets ADD height INTEGER;
UPDATE assets SET size = 0 WHERE size IS NULL;
UPDATE assets SET type = 'application/octet-stream' WHERE type IS NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE assets DROP size;
ALTER TABLE assets DROP type;
ALTER TABLE assets DROP slug;
ALTER TABLE assets DROP width;
ALTER TABLE assets DROP height;

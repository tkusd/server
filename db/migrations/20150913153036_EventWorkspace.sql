-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM events;
DROP TRIGGER IF EXISTS check_action_id ON events;
ALTER TABLE events ADD workspace TEXT NOT NULL DEFAULT '';
ALTER TABLE events DROP COLUMN action_id;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM events;
ALTER TABLE events DROP COLUMN workspace;
ALTER TABLE events ADD action_id UUID NOT NULL REFERENCES actions(id) ON DELETE CASCADE ON UPDATE CASCADE;
CREATE TRIGGER check_action_id
  BEFORE INSERT OR UPDATE
  ON events
  FOR EACH ROW
  EXECUTE PROCEDURE events_check_action_id();
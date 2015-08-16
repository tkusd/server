-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION projects_check_main_screen() RETURNS TRIGGER AS $$
BEGIN
  -- Do nothing if main_screen is same
  IF TG_OP = 'UPDATE' AND NEW.main_screen = OLD.main_screen THEN
    RETURN NEW;
  END IF;

  IF (
    SELECT project_id
    FROM elements
    WHERE id = NEW.main_screen AND element_id IS NULL
    LIMIT 1
  ) = NEW.id THEN
    RETURN NEW;
  END IF;

  RAISE EXCEPTION USING
    errcode = 'foreign_key_violation',
    message = 'Main screen is not owned by the project',
    hint = 'element_not_owned_by_project';

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION events_check_action_id() RETURNS TRIGGER AS $$
BEGIN
  -- Do nothing if action_id is same
  IF TG_OP = 'UPDATE' AND NEW.action_id = OLD.action_id THEN
    RETURN NEW;
  END IF;

  IF (
    SELECT project_id
    FROM actions
    WHERE id = NEW.action_id
    LIMIT 1
  ) = (
    SELECT project_id
    FROM elements
    WHERE id = NEW.element_id
    LIMIT 1
  ) THEN
    RETURN NEW;
  END IF;

  RAISE EXCEPTION USING
    errcode = 'foreign_key_violation',
    message = 'Action is not owned by the project',
    hint = 'action_not_owned_by_project';

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION elements_check_parent_project_id() RETURNS TRIGGER AS $$
BEGIN
  -- Do nothing if element_id is null
  IF NEW.element_id IS NULL THEN
    RETURN NEW;
  END IF;

  -- Do nothing if element_id is same
  IF TG_OP = 'UPDATE' AND NEW.element_id = OLD.element_id THEN
    RETURN NEW;
  END IF;

  IF (
    SELECT project_id
    FROM elements
    WHERE id = NEW.element_id
    LIMIT 1
  ) = NEW.project_id THEN
    RETURN NEW;
  END IF;

  RAISE EXCEPTION USING
    errcode = 'foreign_key_violation',
    message = 'Element is not owned by the project',
    hint = 'element_not_owned_by_project';

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER check_main_screen
  BEFORE INSERT OR UPDATE
  ON projects
  FOR EACH ROW
  EXECUTE PROCEDURE projects_check_main_screen();

CREATE TRIGGER check_action_id
  BEFORE INSERT OR UPDATE
  ON events
  FOR EACH ROW
  EXECUTE PROCEDURE events_check_action_id();

CREATE TRIGGER check_parent_project_id
  BEFORE INSERT OR UPDATE
  ON elements
  FOR EACH ROW
  EXECUTE PROCEDURE elements_check_parent_project_id();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TRIGGER IF EXISTS check_main_screen ON projects;
DROP FUNCTION projects_check_main_screen();

DROP TRIGGER IF EXISTS check_action_id ON events;
DROP FUNCTION events_check_action_id();

DROP TRIGGER IF EXISTS check_parent_project_id ON elements;
DROP FUNCTION elements_check_parent_project_id();

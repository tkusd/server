-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP TRIGGER IF EXISTS check_main_screen ON projects;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION projects_check_main_screen() RETURNS TRIGGER AS $$
BEGIN
  -- Do nothing if main_screen is same
  IF NEW.main_screen IS NULL OR NEW.main_screen = OLD.main_screen THEN
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
-- +goose StatementEnd

CREATE TRIGGER check_main_screen
  BEFORE UPDATE
  ON projects
  FOR EACH ROW
  EXECUTE PROCEDURE projects_check_main_screen();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE elements ALTER COLUMN index SET DEFAULT 1;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION elements_update_index() RETURNS TRIGGER AS $$
BEGIN
  SELECT COALESCE((
    SELECT index
    FROM elements
    WHERE element_id = NEW.element_id AND project_id = NEW.project_id
    ORDER BY index desc
    LIMIT 1
  ) + 1, 1) INTO NEW.index;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION elements_replace_main_screen() RETURNS TRIGGER AS $$
BEGIN
  IF OLD.element_id IS NOT NULL THEN
    RETURN OLD;
  END IF;

  IF (SELECT (main_screen = OLD.id) FROM projects WHERE id = OLD.project_id) THEN
    UPDATE projects
    SET main_screen = (
      SELECT id
      FROM elements
      WHERE project_id = OLD.project_id
        AND element_id IS NULL
        AND id <> OLD.id
      ORDER BY index
      LIMIT 1
    )
    WHERE id = OLD.project_id;
  END IF;

  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION elements_create_main_screen() RETURNS TRIGGER AS $$
BEGIN
  IF NEW.element_id IS NOT NULL THEN
    RETURN NEW;
  END IF;

  IF (SELECT main_screen IS NULL FROM projects WHERE id = NEW.project_id) THEN
    UPDATE projects SET main_screen = NEW.id WHERE id = NEW.project_id;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER update_index
  BEFORE INSERT
  ON elements
  FOR EACH ROW
  EXECUTE PROCEDURE elements_update_index();

CREATE TRIGGER replace_main_screen
  BEFORE DELETE
  ON elements
  FOR EACH ROW
  EXECUTE PROCEDURE elements_replace_main_screen();

CREATE TRIGGER create_main_screen
  AFTER INSERT
  ON elements
  FOR EACH ROW
  EXECUTE PROCEDURE elements_create_main_screen();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TRIGGER IF EXISTS update_index ON elements;
DROP FUNCTION elements_update_index();

DROP TRIGGER IF EXISTS replace_main_screen ON elements;
DROP FUNCTION elements_replace_main_screen();

DROP TRIGGER IF EXISTS create_main_screen ON elements;
DROP FUNCTION elements_create_main_screen();

ALTER TABLE elements ALTER COLUMN index DROP DEFAULT;

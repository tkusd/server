-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION elements_update_index() RETURNS TRIGGER AS $$
BEGIN
  IF NEW.element_id IS NULL THEN
    SELECT COALESCE((
      SELECT index
      FROM elements
      WHERE element_id IS NULL AND project_id = NEW.project_id
      ORDER BY index desc
      LIMIT 1
    ) + 1, 1) INTO NEW.index;
  ELSE
    SELECT COALESCE((
      SELECT index
      FROM elements
      WHERE element_id = NEW.element_id AND project_id = NEW.project_id
      ORDER BY index desc
      LIMIT 1
    ) + 1, 1) INTO NEW.index;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

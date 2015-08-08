-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION elements_update_project_date() RETURNS TRIGGER AS $$
DECLARE
	rec elements;
BEGIN
	IF TG_OP = 'DELETE' THEN
		rec := OLD;
	ELSE
		rec := NEW;
	END IF;

	UPDATE projects SET updated_at = current_timestamp WHERE id = rec.project_id;
	RETURN rec;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER update_project
	AFTER INSERT OR UPDATE OR DELETE
	ON elements
	FOR EACH ROW
	EXECUTE PROCEDURE elements_update_project_date();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TRIGGER IF EXISTS update_project ON elements;
DROP FUNCTION elements_update_project_date();

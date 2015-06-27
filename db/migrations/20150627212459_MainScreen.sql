-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE projects ADD main_screen UUID REFERENCES elements(id) ON DELETE SET NULL ON UPDATE CASCADE;
UPDATE projects
  SET main_screen = (SELECT id FROM elements WHERE project_id = projects.id AND element_id IS NULL ORDER BY order_id LIMIT 1)
  WHERE main_screen IS NULL
  AND (SELECT 1 FROM elements WHERE project_id = projects.id LIMIT 1) > 0;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE projects DROP main_screen;
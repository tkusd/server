-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE OR REPLACE VIEW public_users AS
	SELECT id, name, avatar
	FROM users;

CREATE OR REPLACE VIEW projects_extended AS
	SELECT
		projects.*,
		users.name AS user_name,
		users.avatar AS user_avatar,
		(SELECT count(id) FROM elements WHERE project_id = projects.id) AS element_count,
		(SELECT count(id) FROM assets WHERE project_id = projects.id) AS asset_count
	FROM projects
	JOIN users
	ON users.id = projects.user_id;

CREATE OR REPLACE VIEW users_extended AS
	SELECT
		users.*,
		(SELECT count(id) FROM projects WHERE user_id = users.id) AS total_project_count,
		(SELECT count(id) FROM projects WHERE user_id = users.id AND is_private = FALSE) as public_project_count
	FROM users;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP VIEW users_extended;
DROP VIEW projects_extended;
DROP VIEW public_users;

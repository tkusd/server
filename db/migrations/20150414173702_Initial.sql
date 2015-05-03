-- +goose Up
-- SQL in section Up is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
	id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(100) NOT NULL DEFAULT '',
	password CHAR(60) NOT NULL DEFAULT '',
	email VARCHAR(254) NOT NULL UNIQUE,
	avatar VARCHAR(254) NOT NULL DEFAULT '',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_activated BOOLEAN NOT NULL DEFAULT FALSE,
	activation_token CHAR(64)
);

CREATE TABLE IF NOT EXISTS reset_tokens (
	id CHAR(64) NOT NULL PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tokens (
	id CHAR(64) NOT NULL PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
	id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	title VARCHAR(255) NOT NULL,
	description TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_private BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS project_members (
	id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	role BIT(3) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS elements (
	id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
	project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	element_id UUID REFERENCES elements(id) ON DELETE CASCADE,
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	attributes JSONB NOT NULL DEFAULT '{}'
);

-- +goose Down
-- SQL section Down is executed when this migration is rolled back
DROP TABLE IF EXISTS elements;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS reset_tokens;
DROP TABLE IF EXISTS users;
-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TYPE user_role AS ENUM ('admin', 'member', 'viewer');

CREATE TABLE users (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT        NOT NULL,
    email       TEXT        NOT NULL UNIQUE,
    role        user_role   NOT NULL DEFAULT 'member',
    locale      TEXT        NOT NULL DEFAULT 'en',
    preferences JSONB       NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE project_status AS ENUM ('active', 'archived');

CREATE TABLE projects (
    id          UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT           NOT NULL,
    description TEXT           NOT NULL DEFAULT '',
    status      project_status NOT NULL DEFAULT 'active',
    settings    JSONB          NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE project_members (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID        NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    role       user_role   NOT NULL DEFAULT 'member',
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, user_id)
);

CREATE INDEX idx_project_members_project_id ON project_members(project_id);
CREATE INDEX idx_project_members_user_id    ON project_members(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
DROP TYPE  IF EXISTS project_status;
DROP TYPE  IF EXISTS user_role;
DROP EXTENSION IF EXISTS vector;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd

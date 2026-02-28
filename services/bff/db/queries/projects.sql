-- name: GetProject :one
SELECT * FROM projects WHERE id = $1 LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProjects :one
SELECT COUNT(*) FROM projects;

-- name: CreateProject :one
INSERT INTO projects (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET
    name        = COALESCE(NULLIF($2, ''), name),
    description = COALESCE(NULLIF($3, ''), description),
    status      = COALESCE($4, status),
    updated_at  = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: ListProjectMembers :many
SELECT pm.*, u.name AS user_name, u.email AS user_email
FROM project_members pm
JOIN users u ON u.id = pm.user_id
WHERE pm.project_id = $1;

-- name: AddProjectMember :one
INSERT INTO project_members (project_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (project_id, user_id) DO UPDATE SET role = EXCLUDED.role
RETURNING *;

-- name: RemoveProjectMember :exec
DELETE FROM project_members
WHERE project_id = $1 AND user_id = $2;

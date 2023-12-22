-- name: CreateTheme :one
INSERT INTO themes (name, subsection_id, author_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTheme :one
SELECT id, name, subsection_id, author_id
FROM themes
WHERE id = $1;

-- name: GetThemes :many
SELECT id, name, subsection_id, author_id
FROM themes
WHERE subsection_id = $1;

-- name: UpdateTheme :one
UPDATE themes
SET name = $2, subsection_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTheme :exec
DELETE FROM themes
WHERE id = $1;
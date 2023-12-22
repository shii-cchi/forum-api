-- name: CreateSection :one
INSERT INTO sections (name)
VALUES ($1)
RETURNING *;

-- name: GetSection :one
SELECT id, name
FROM sections
WHERE id = $1;

-- name: GetSections :many
SELECT id, name
FROM sections;

-- name: UpdateSectionName :one
UPDATE sections
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSection :exec
DELETE FROM sections
WHERE id = $1;
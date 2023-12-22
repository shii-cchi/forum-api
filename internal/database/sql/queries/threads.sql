-- name: CreateThread :one
INSERT INTO threads (name, theme_id, author_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetThread :one
SELECT id, name, theme_id, author_id
FROM threads
WHERE id = $1;

-- name: GetThreads :many
SELECT id, name, theme_id, author_id
FROM threads
WHERE theme_id = $1;

-- name: UpdateThread :one
UPDATE threads
SET name = $2, theme_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteThread :exec
DELETE FROM threads
WHERE id = $1;
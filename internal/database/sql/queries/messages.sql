-- name: CreateMessage :one
INSERT INTO messages (name, thread_id, author_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMessage :one
SELECT id, name, thread_id, author_id
FROM messages
WHERE id = $1;

-- name: GetMessages :many
SELECT id, name, thread_id, author_id
FROM messages
WHERE thread_id = $1;

-- name: UpdateMessage :one
UPDATE messages
SET name = $2, thread_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;
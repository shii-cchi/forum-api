-- name: CreateSubsection :one
INSERT INTO subsections (name, section_id, author_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSubsection :one
SELECT id, name, section_id, author_id
FROM subsections
WHERE id = $1;

-- name: GetSubsections :many
SELECT id, name, section_id, author_id
FROM subsections
WHERE section_id = $1;

-- name: UpdateSubsection :one
UPDATE subsections
SET name = $2, section_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSubsection :exec
DELETE FROM subsections
WHERE id = $1;
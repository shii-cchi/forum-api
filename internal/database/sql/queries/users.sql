-- name: CreateUser :one
INSERT INTO users (email, password, login)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddToken :exec
UPDATE users
SET token = $2
WHERE id = $1;

-- name: CheckUserIsExist :one
SELECT COUNT(*)
FROM users
WHERE email = $1 OR login = $2;
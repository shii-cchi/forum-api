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

-- name: LogoutUser :exec
UPDATE users
SET token = ''
WHERE id = $1;

-- name: CheckDataToLogin :one
SELECT *
FROM users
WHERE (email = $1 AND password = $2) OR (login = $3 AND password = $2);
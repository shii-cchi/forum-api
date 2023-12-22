-- name: CreateUser :one
INSERT INTO users (email, password, login, role_id)
VALUES ($1, $2, $3, 2)
RETURNING *;

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
WHERE email = $1 OR login = $2;

-- name: FindUserById :one
SELECT *
FROM users
WHERE id = $1;
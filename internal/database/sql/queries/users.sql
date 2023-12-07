-- name: CreateUser :one
INSERT INTO users (email, password, login, role_id)
VALUES ($1, $2, $3, 2)
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
WHERE email = $1 OR login = $2;

-- name: FindUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetRole :one
SELECT roles.name
FROM users
JOIN roles ON users.role_id = roles.id
WHERE users.id = $1;

-- name: GetPermissions :many
SELECT permissions.name
FROM roles
JOIN roles_permissions ON roles_permissions.role_id = roles.id
JOIN permissions ON roles_permissions.permission_id = permissions.id
WHERE roles.name = $1;
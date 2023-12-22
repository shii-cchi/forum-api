-- name: GetRole :one
SELECT roles.name
FROM users
         JOIN roles ON users.role_id = roles.id
WHERE users.id = $1;
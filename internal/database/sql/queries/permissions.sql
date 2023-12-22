-- name: GetPermissions :many
SELECT permissions.name
FROM roles
         JOIN roles_permissions ON roles_permissions.role_id = roles.id
         JOIN permissions ON roles_permissions.permission_id = permissions.id
WHERE roles.name = $1;
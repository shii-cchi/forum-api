-- +goose Up

INSERT INTO roles (id, name)
VALUES (1, 'admin'),
       (2, 'user');

INSERT INTO permissions (id, name)
VALUES (1, 'TEST1'),
       (2, 'TEST2'),
       (3, 'TEST3');

INSERT INTO roles_permissions (role_id, permission_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (2, 1),
       (2, 2);

-- +goose Down


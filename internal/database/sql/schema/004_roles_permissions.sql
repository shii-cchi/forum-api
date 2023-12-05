-- +goose Up

CREATE TABLE roles_permissions (
    role_id bigint NOT NULL,
    permission_id bigint NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);

-- +goose Down

DROP TABLE roles_permissions;
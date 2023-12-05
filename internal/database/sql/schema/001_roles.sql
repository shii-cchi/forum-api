-- +goose Up

CREATE TABLE roles (
    id bigint NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down

DROP TABLE roles;
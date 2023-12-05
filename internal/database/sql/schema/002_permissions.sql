-- +goose Up

CREATE TABLE permissions (
    id bigint NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down

DROP TABLE permissions;
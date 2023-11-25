-- +goose Up

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    login TEXT NOT NULL,
    token VARCHAR(64) UNIQUE NOT NULL
);

-- +goose Down

DROP TABLE users;
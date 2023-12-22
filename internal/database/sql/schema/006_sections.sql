-- +goose Up

CREATE TABLE sections (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down

DROP TABLE sections;
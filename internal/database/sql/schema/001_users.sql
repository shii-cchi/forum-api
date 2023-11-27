-- +goose Up

DROP EXTENSION IF EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    login TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL DEFAULT ''
);

-- +goose Down

DROP TABLE users;
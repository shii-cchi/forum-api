-- +goose Up

CREATE TABLE threads (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    theme_id uuid NOT NULL,
    author_id uuid NOT NULL,
    FOREIGN KEY (theme_id) REFERENCES themes(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE threads;
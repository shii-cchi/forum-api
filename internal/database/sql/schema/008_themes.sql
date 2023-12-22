-- +goose Up

CREATE TABLE themes (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    subsection_id uuid NOT NULL,
    author_id uuid NOT NULL,
    FOREIGN KEY (subsection_id) REFERENCES subsections(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE themes;
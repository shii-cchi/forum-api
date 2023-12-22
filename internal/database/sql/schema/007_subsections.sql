-- +goose Up

CREATE TABLE subsections (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    section_id uuid NOT NULL,
    author_id uuid NOT NULL,
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE subsections;
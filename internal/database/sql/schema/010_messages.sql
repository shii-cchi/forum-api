-- +goose Up

CREATE TABLE messages (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    thread_id uuid NOT NULL,
    author_id uuid NOT NULL,
    FOREIGN KEY (thread_id) REFERENCES thread(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE messages;
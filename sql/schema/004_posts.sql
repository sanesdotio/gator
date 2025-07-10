-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    title text NOT NULL,
    url text NOT NULL,
    description text,
    published_at timestamp,
    feed_id uuid NOT NULL REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
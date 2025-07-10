-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    last_fetched_at timestamp,
    name text UNIQUE NOT NULL,
    url text UNIQUE NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
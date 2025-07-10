-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;
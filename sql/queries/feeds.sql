-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, $2, $3, $4, $5, $6, NULL)
RETURNING *;

-- name: GetFeedByURL :one
SELECT feeds.id, feeds.name FROM feeds WHERE url = $1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedOwner :one
SELECT users.name
FROM users
JOIN feeds ON feeds.user_id = users.id
WHERE feeds.id = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $3
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, url, last_fetched_at, updated_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
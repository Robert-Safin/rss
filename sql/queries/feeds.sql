-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: FindFeedUser :one
SELECT * FROM users where id = $1;

-- name: GetAllFeeds :many
SELECT *, users.name as user_name FROM feeds JOIN users ON feeds.user_id = users.id;

-- name: FindFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: CreateUser :one
INSERT INTO users(username, email, password, created, updated)
VALUES ($1, $2, $3, $4, $5)
	RETURN id, username, email, created, updated;

-- name: GetUser :one
SELECT id, username, email, created, updated
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, email, created, updated
FROM users
ORDER BY id;

-- name: CreateBlog :one
INSERT INTO blogs(title, content, user_id, created, updated)
VALUES ($1, $2, $3, $4, $5)
	RETURNING id, title, user_id, created, updated;

-- name: ListBlogs :many
SELECT id, title, content, user_id, created, updated
FROM blogs
ORDER BY id;


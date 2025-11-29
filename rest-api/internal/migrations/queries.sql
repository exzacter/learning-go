-- name: CreateUser :one
INSERT INTO users(username, email, password, created, updated)
VALUES ($1, $2, $3, $4, $5)
	RETURNING id, username, email, created, updated;

-- name: GetUser :one
SELECT id, username, email, created, updated
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, email, created, updated
FROM users
ORDER BY id;

-- name: GetUserByUsernameOrEmail :one
SELECT id, username, email, created, updated, password
FROM users
WHERE username = $1 OR email = $1;

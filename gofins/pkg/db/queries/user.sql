-- name: GetUser :one
SELECT id, name, created_at, is_admin
FROM users
WHERE name = $1;

-- name: GetUserByID :one
SELECT id, name, created_at, is_admin
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, created_at, is_admin
FROM users
ORDER BY created_at ASC;

-- name: CreateUser :one
INSERT INTO users (id, name, created_at, is_admin)
VALUES ($1, $2, NOW(), $3)
RETURNING id, name, created_at, is_admin;

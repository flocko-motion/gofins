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

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUserAdmin :one
UPDATE users
SET is_admin = $1
WHERE id = $2
RETURNING id, name, created_at, is_admin;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteUserRatings :exec
DELETE FROM user_ratings WHERE user_id = $1;

-- name: DeleteUserFavorites :exec
DELETE FROM user_favorites WHERE user_id = $1;

-- name: IsFavorite :one
SELECT EXISTS(SELECT 1 FROM user_favorites WHERE user_id = $1 AND ticker = $2);

-- name: AddFavorite :exec
INSERT INTO user_favorites (user_id, ticker) VALUES ($1, $2);

-- name: RemoveFavorite :exec
DELETE FROM user_favorites WHERE user_id = $1 AND ticker = $2;

-- name: GetFavorites :many
SELECT ticker FROM user_favorites WHERE user_id = $1 ORDER BY created_at DESC;

-- name: AddRating :one
INSERT INTO user_ratings (user_id, ticker, rating, notes)
VALUES ($1, $2, $3, $4)
RETURNING id, ticker, rating, notes, created_at;

-- name: GetLatestRating :one
SELECT id, ticker, rating, notes, created_at
FROM user_ratings
WHERE user_id = $1 AND ticker = $2
ORDER BY created_at DESC
LIMIT 1;

-- name: GetRatingHistory :many
SELECT id, ticker, rating, notes, created_at
FROM user_ratings
WHERE user_id = $1 AND ticker = $2
ORDER BY created_at DESC;

-- name: GetAllLatestRatings :many
SELECT DISTINCT ON (ticker) id, ticker, rating, notes, created_at
FROM user_ratings
WHERE user_id = $1
ORDER BY ticker, created_at DESC;

-- name: DeleteRating :exec
DELETE FROM user_ratings WHERE user_id = $1 AND id = $2;

-- name: GetAllNotesChronological :many
SELECT id, ticker, rating, notes, created_at
FROM user_ratings
WHERE user_id = $1 AND notes IS NOT NULL AND notes != ''
ORDER BY created_at DESC;

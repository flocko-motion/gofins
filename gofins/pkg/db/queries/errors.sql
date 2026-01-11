-- name: InsertError :exec
INSERT INTO errors (source, error_type, message, details)
VALUES ($1, $2, $3, $4);

-- name: GetRecentErrors :many
SELECT id, timestamp, source, error_type, message, details
FROM errors
ORDER BY timestamp DESC
LIMIT $1;

-- name: GetErrorsBySource :many
SELECT id, timestamp, source, error_type, message, details
FROM errors
WHERE source = $1
ORDER BY timestamp DESC
LIMIT $2;

-- name: CountErrorsSince :one
SELECT COUNT(*) FROM errors WHERE timestamp >= $1;

-- name: ClearOldErrors :execrows
DELETE FROM errors WHERE timestamp < $1;

-- name: GetErrorByID :one
SELECT id, timestamp, source, error_type, message, details
FROM errors
WHERE id = $1;

-- name: ClearAllErrors :execrows
DELETE FROM errors;

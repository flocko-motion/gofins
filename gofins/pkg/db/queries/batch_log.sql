-- name: StartBatchUpdate :one
INSERT INTO batch_update_log (updater_name, started_at, status)
VALUES ($1, $2, 'running')
RETURNING id;

-- name: CompleteBatchUpdate :exec
UPDATE batch_update_log
SET completed_at = $1,
    status = 'completed',
    symbols_processed = $2,
    symbols_updated = $3
WHERE id = $4;

-- name: FailBatchUpdate :exec
UPDATE batch_update_log
SET completed_at = $1,
    status = 'failed',
    error_message = $2
WHERE id = $3;

-- name: GetLastBatchUpdate :one
SELECT id, updater_name, started_at, completed_at, status, 
       symbols_processed, symbols_updated, error_message
FROM batch_update_log
WHERE updater_name = $1
  AND status = 'completed'
ORDER BY started_at DESC
LIMIT 1;

-- name: DeleteBatchUpdate :exec
DELETE FROM batch_update_log
WHERE updater_name = $1;

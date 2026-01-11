-- name: CreateAnalysisPackage :exec
INSERT INTO analysis_packages (
    id, name, created_at, interval, time_from, time_to,
    hist_bins, hist_min, hist_max, mcap_min, inception_max, status, user_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: UpdateAnalysisPackageStatus :exec
UPDATE analysis_packages 
SET status = $1, symbol_count = $2 
WHERE id = $3 AND user_id = $4;

-- name: VerifyPackageOwnership :one
SELECT EXISTS(SELECT 1 FROM analysis_packages WHERE id = $1 AND user_id = $2);

-- name: SaveAnalysisResult :exec
INSERT INTO analysis_results (
    package_id, ticker, count, mean, stddev, variance, min, max, histogram
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetAnalysisResults :many
SELECT ar.package_id, ar.ticker, ar.count, ar.mean, ar.stddev, ar.variance, ar.min, ar.max, s.inception
FROM analysis_results ar
JOIN symbols s ON ar.ticker = s.ticker
WHERE ar.package_id = $1
ORDER BY ar.mean DESC;

-- name: GetAnalysisPackage :one
SELECT id, name, created_at, interval, time_from, time_to,
       hist_bins, hist_min, hist_max, mcap_min, inception_max, symbol_count, status, user_id
FROM analysis_packages
WHERE id = $1 AND user_id = $2;

-- name: ListAnalysisPackages :many
SELECT id, name, created_at, interval, time_from, time_to,
       hist_bins, hist_min, hist_max, mcap_min, inception_max, symbol_count, status, user_id
FROM analysis_packages
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateAnalysisPackageName :exec
UPDATE analysis_packages 
SET name = $1 
WHERE id = $2 AND user_id = $3;

-- name: DeleteAnalysisPackage :exec
DELETE FROM analysis_packages 
WHERE id = $1 AND user_id = $2;

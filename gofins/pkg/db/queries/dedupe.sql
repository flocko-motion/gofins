-- name: GetSymbolsWithCIK :many
SELECT ticker, exchange, cik, primary_listing
FROM symbols
WHERE cik IS NOT NULL AND cik != ''
  AND (type = $1 OR type IS NULL)
ORDER BY cik, exchange;

-- name: GetStockSymbolsWithoutCIK :many
SELECT ticker, name, oldest_price, primary_listing
FROM symbols
WHERE (cik IS NULL OR cik = '')
  AND (type = $1 OR type IS NULL)
  AND name IS NOT NULL
ORDER BY name, oldest_price ASC NULLS LAST;

-- name: GetStockSymbolsForNameDedupe :many
SELECT ticker, name, oldest_price, primary_listing
FROM symbols
WHERE (type = $1 OR type IS NULL)
  AND name IS NOT NULL
ORDER BY name, oldest_price ASC NULLS LAST;

-- name: ResetPrimaryListings :execrows
UPDATE symbols SET primary_listing = NULL;

-- name: SetPrimaryListing :exec
UPDATE symbols SET primary_listing = '' WHERE ticker = $1;

-- name: SetSecondaryListings :exec
UPDATE symbols SET primary_listing = $1 WHERE ticker = ANY($2::text[]);

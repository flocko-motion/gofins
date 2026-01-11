-- name: GetSymbol :one
SELECT s.ticker, s.exchange, s.last_price_update, s.last_profile_update,
       s.last_price_status, s.last_profile_status,
       s.name, s.type, s.currency, s.sector, s.industry, s.country,
       s.description, s.website, s.isin, s.cik, s.inception, s.oldest_price, 
       s.is_actively_trading, s.market_cap, s.primary_listing,
       s.ath12m, s.current_price_usd, s.current_price_time,
       COALESCE(f.ticker IS NOT NULL, false) as is_favorite
FROM symbols s
LEFT JOIN user_favorites f ON s.ticker = f.ticker
WHERE s.ticker = $1;

-- name: GetAllTickers :many
SELECT ticker FROM symbols;

-- name: GetStaleProfiles :many
SELECT ticker FROM symbols
WHERE (last_profile_update IS NULL OR last_profile_update < $1)
  AND (type IS NULL OR type != $2)
  AND (primary_listing IS NULL OR primary_listing = '')
ORDER BY last_profile_update ASC NULLS FIRST
LIMIT $3;

-- name: CountSymbols :one
SELECT COUNT(*) FROM symbols;

-- name: CountActivelyTrading :one
SELECT COUNT(*) FROM symbols WHERE is_actively_trading = true;

-- name: CountStaleProfiles :one
SELECT COUNT(*) FROM symbols
WHERE (last_profile_update IS NULL OR last_profile_update < $1)
  AND (type IS NULL OR type != $2);

-- name: GetOldestProfileUpdate :one
SELECT MIN(last_profile_update) FROM symbols
WHERE last_profile_update IS NOT NULL;

-- name: ResetQuoteTimestamps :execrows
UPDATE symbols 
SET current_price_time = NULL;

-- name: ResetPriceTimestamps :execrows
UPDATE symbols 
SET last_price_update = NULL, last_price_status = NULL;

-- name: ResetProfileTimestamps :execrows
UPDATE symbols 
SET last_profile_update = NULL, last_profile_status = NULL;

-- name: ResetSymbol :execrows
UPDATE symbols 
SET last_price_update = NULL, last_price_status = NULL,
    last_profile_update = NULL, last_profile_status = NULL
WHERE ticker = $1;

-- name: ResetIndexTimestamps :execrows
UPDATE symbols 
SET last_price_update = NULL, 
    last_price_status = NULL,
    last_profile_update = NULL,
    last_profile_status = NULL,
    is_actively_trading = true
WHERE type = $1;

-- name: GetAllSymbolCurrencies :many
SELECT ticker, COALESCE(currency, 'USD') as currency
FROM symbols;

-- name: MarkStaleProfilesAsNotFound :execrows
UPDATE symbols
SET last_profile_status = $1
WHERE last_profile_update < $2
   OR last_profile_update IS NULL;

-- name: GetTickersNeedingProfileUpdate :many
SELECT ticker 
FROM symbols 
WHERE last_profile_update IS NULL 
   OR last_profile_update < $1;

-- name: GetTickersNeedingQuoteUpdate :many
SELECT ticker 
FROM symbols 
WHERE current_price_time IS NULL 
   OR current_price_time < $1;

-- name: GetOldestPriceDate :one
SELECT MIN(date) 
FROM monthly_prices 
WHERE symbol_ticker = $1;

-- name: GetFilteredTickers :many
SELECT DISTINCT s.ticker FROM symbols s
WHERE s.is_actively_trading = true
  AND s.type = ANY($1::text[])
  AND ($2::BIGINT IS NULL OR s.market_cap >= $2)
  AND ($3::TIMESTAMP IS NULL OR s.inception <= $3)
  AND s.last_price_status = $4
  AND s.last_price_update IS NOT NULL
  AND s.exchange NOT IN ('OTC','PINK', 'GREY', 'OTCQB', 'OTCQX')
ORDER BY s.ticker;

-- name: GetTickersWithPrices :many
SELECT DISTINCT symbol_ticker FROM monthly_prices 
ORDER BY symbol_ticker 
LIMIT $1;

-- name: GetSymbolsWithStalePrices :many
SELECT ticker, currency FROM symbols
WHERE (last_price_update IS NULL OR last_price_update < $1)
  AND is_actively_trading = true
  AND type = ANY($2::text[])
ORDER BY last_price_update ASC NULLS FIRST
LIMIT $3;

-- name: CountStalePrices :one
SELECT COUNT(*) FROM symbols
WHERE (last_price_update IS NULL OR last_price_update < $1)
  AND is_actively_trading = true
  AND type = ANY($2::text[]);

-- name: GetOldestPriceUpdate :one
SELECT MIN(last_price_update) FROM symbols
WHERE last_price_update IS NOT NULL
  AND is_actively_trading = true;

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/f"
	"github.com/flocko-motion/gofins/pkg/types"
	"github.com/lib/pq"
)

// GetOldestPriceDate returns the oldest price date for a ticker from monthly_prices
func GetOldestPriceDate(ctx context.Context, ticker string) (*time.Time, error) {
	result, err := genQ().GetOldestPriceDate(ctx, ticker)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest price date: %w", err)
	}

	// Handle interface{} return from sqlc
	if result == nil {
		return nil, nil
	}
	if t, ok := result.(time.Time); ok {
		return &t, nil
	}
	return nil, nil
}

// AppendSinglePrice adds a single price point to the price history
// Used for incremental updates when we just need to add the latest period
func AppendSinglePrice(price types.PriceData, interval types.PriceInterval) error {
	db := Db()
	tableName := string(interval) + "_prices"

	query := fmt.Sprintf(`
		INSERT INTO %s (symbol_ticker, date, open, high, low, avg, close, yoy, open_orig, high_orig, low_orig, avg_orig, close_orig)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (symbol_ticker, date) DO UPDATE SET
			open = EXCLUDED.open,
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			avg = EXCLUDED.avg,
			close = EXCLUDED.close,
			yoy = EXCLUDED.yoy,
			open_orig = EXCLUDED.open_orig,
			high_orig = EXCLUDED.high_orig,
			low_orig = EXCLUDED.low_orig,
			avg_orig = EXCLUDED.avg_orig,
			close_orig = EXCLUDED.close_orig
	`, tableName)

	_, err := db.conn.Exec(query,
		price.SymbolTicker, price.Date, price.Open, price.High, price.Low,
		price.Avg, price.Close, price.YoY, price.OpenOrig, price.HighOrig,
		price.LowOrig, price.AvgOrig, price.CloseOrig)

	return err
}

// PutMonthlyPrices batch inserts monthly price data using bulk INSERT
func PutMonthlyPrices(prices []types.PriceData) error {
	db := Db()
	if len(prices) == 0 {
		return nil
	}

	// Build bulk INSERT with all VALUES in one query
	// Split into chunks of 1000 to avoid parameter limits
	chunkSize := 1000
	for i := 0; i < len(prices); i += chunkSize {
		end := i + chunkSize
		if end > len(prices) {
			end = len(prices)
		}
		chunk := prices[i:end]

		// Build VALUES list: ($1,$2,...), ($13,$14,...), ...
		valueStrings := make([]string, 0, len(chunk))
		valueArgs := make([]interface{}, 0, len(chunk)*13)

		for idx, p := range chunk {
			paramOffset := idx * 13
			valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				paramOffset+1, paramOffset+2, paramOffset+3, paramOffset+4,
				paramOffset+5, paramOffset+6, paramOffset+7, paramOffset+8,
				paramOffset+9, paramOffset+10, paramOffset+11, paramOffset+12, paramOffset+13))
			valueArgs = append(valueArgs, p.Date, p.Open, p.High, p.Low, p.Avg, p.Close, p.YoY, p.SymbolTicker,
				p.OpenOrig, p.HighOrig, p.LowOrig, p.AvgOrig, p.CloseOrig)
		}

		query := fmt.Sprintf(`
			INSERT INTO monthly_prices (date, open, high, low, avg, close, yoy, symbol_ticker, open_orig, high_orig, low_orig, avg_orig, close_orig)
			VALUES %s
			ON CONFLICT (date, symbol_ticker) DO UPDATE SET
				open = EXCLUDED.open,
				high = EXCLUDED.high,
				low = EXCLUDED.low,
				avg = EXCLUDED.avg,
				close = EXCLUDED.close,
				yoy = EXCLUDED.yoy,
				open_orig = EXCLUDED.open_orig,
				high_orig = EXCLUDED.high_orig,
				low_orig = EXCLUDED.low_orig,
				avg_orig = EXCLUDED.avg_orig,
				close_orig = EXCLUDED.close_orig
		`, joinStrings(valueStrings, ","))

		_, err := db.conn.Exec(query, valueArgs...)
		if err != nil {
			return fmt.Errorf("failed to batch insert monthly prices: %w", err)
		}
	}

	return nil
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// PutWeeklyPrices batch inserts weekly price data using bulk INSERT
func PutWeeklyPrices(prices []types.PriceData) error {
	db := Db()
	if len(prices) == 0 {
		return nil
	}

	// Build bulk INSERT with all VALUES in one query
	// Split into chunks of 1000 to avoid parameter limits
	chunkSize := 1000
	for i := 0; i < len(prices); i += chunkSize {
		end := i + chunkSize
		if end > len(prices) {
			end = len(prices)
		}
		chunk := prices[i:end]

		// Build VALUES list: ($1,$2,...), ($13,$14,...), ...
		valueStrings := make([]string, 0, len(chunk))
		valueArgs := make([]interface{}, 0, len(chunk)*13)

		for idx, p := range chunk {
			paramOffset := idx * 13
			valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				paramOffset+1, paramOffset+2, paramOffset+3, paramOffset+4,
				paramOffset+5, paramOffset+6, paramOffset+7, paramOffset+8,
				paramOffset+9, paramOffset+10, paramOffset+11, paramOffset+12, paramOffset+13))
			valueArgs = append(valueArgs, p.Date, p.Open, p.High, p.Low, p.Avg, p.Close, p.YoY, p.SymbolTicker,
				p.OpenOrig, p.HighOrig, p.LowOrig, p.AvgOrig, p.CloseOrig)
		}

		query := fmt.Sprintf(`
			INSERT INTO weekly_prices (date, open, high, low, avg, close, yoy, symbol_ticker, open_orig, high_orig, low_orig, avg_orig, close_orig)
			VALUES %s
			ON CONFLICT (date, symbol_ticker) DO UPDATE SET
				open = EXCLUDED.open,
				high = EXCLUDED.high,
				low = EXCLUDED.low,
				avg = EXCLUDED.avg,
				close = EXCLUDED.close,
				yoy = EXCLUDED.yoy,
				open_orig = EXCLUDED.open_orig,
				high_orig = EXCLUDED.high_orig,
				low_orig = EXCLUDED.low_orig,
				avg_orig = EXCLUDED.avg_orig,
				close_orig = EXCLUDED.close_orig
		`, joinStrings(valueStrings, ","))

		_, err := db.conn.Exec(query, valueArgs...)
		if err != nil {
			return fmt.Errorf("failed to batch insert weekly prices: %w", err)
		}
	}

	return nil
}

// GetPrices retrieves price data for a symbol at the specified interval
func GetPrices(ticker string, from, to time.Time, interval types.PriceInterval) ([]types.PriceData, error) {
	db := Db()
	tableName := string(interval) + "_prices"

	query := fmt.Sprintf(`
		SELECT date, open, high, low, avg, close, yoy, symbol_ticker, open_orig, high_orig, low_orig, avg_orig, close_orig
		FROM %s
		WHERE symbol_ticker = $1 AND date >= $2 AND date <= $3
		ORDER BY date ASC
	`, tableName)

	rows, err := db.conn.Query(query, ticker, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []types.PriceData
	for rows.Next() {
		var p types.PriceData
		if err := rows.Scan(&p.Date, &p.Open, &p.High, &p.Low, &p.Avg, &p.Close, &p.YoY, &p.SymbolTicker,
			&p.OpenOrig, &p.HighOrig, &p.LowOrig, &p.AvgOrig, &p.CloseOrig); err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}

	return prices, rows.Err()
}

// GetMonthlyPrices retrieves monthly prices for a symbol
func GetMonthlyPrices(ticker string, from, to time.Time) ([]types.PriceData, error) {
	return GetPrices(ticker, from, to, types.IntervalMonthly)
}

// GetWeeklyPrices retrieves weekly prices for a symbol
func GetWeeklyPrices(ticker string, from, to time.Time) ([]types.PriceData, error) {
	return GetPrices(ticker, from, to, types.IntervalWeekly)
}

// GetLatestPriceDate returns the most recent price date for a symbol at the specified interval
// Returns nil if no prices exist
func GetLatestPriceDate(ticker string, interval types.PriceInterval) (*time.Time, error) {
	db := Db()
	tableName := string(interval) + "_prices"

	query := fmt.Sprintf(`
		SELECT MAX(date) 
		FROM %s
		WHERE symbol_ticker = $1
	`, tableName)

	var latestDate *time.Time
	err := db.conn.QueryRow(query, ticker).Scan(&latestDate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return latestDate, nil
}

// GetPricesBatch retrieves price data for multiple symbols in a single query
// Returns a map of ticker -> []PriceData
func GetPricesBatch(tickers []string, from, to time.Time, interval types.PriceInterval) (map[string][]types.PriceData, error) {
	db := Db()
	if len(tickers) == 0 {
		return make(map[string][]types.PriceData), nil
	}

	tableName := string(interval) + "_prices"

	query := fmt.Sprintf(`
		SELECT date, open, high, low, avg, close, yoy, symbol_ticker, open_orig, high_orig, low_orig, avg_orig, close_orig
		FROM %s
		WHERE symbol_ticker = ANY($1) AND date >= $2 AND date <= $3
		ORDER BY symbol_ticker, date ASC
	`, tableName)

	rows, err := db.conn.Query(query, pq.Array(tickers), from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]types.PriceData)
	for rows.Next() {
		var p types.PriceData
		if err := rows.Scan(&p.Date, &p.Open, &p.High, &p.Low, &p.Avg, &p.Close, &p.YoY, &p.SymbolTicker,
			&p.OpenOrig, &p.HighOrig, &p.LowOrig, &p.AvgOrig, &p.CloseOrig); err != nil {
			return nil, err
		}
		result[p.SymbolTicker] = append(result[p.SymbolTicker], p)
	}

	return result, rows.Err()
}

// GetFilteredTickers returns tickers matching filters (for analysis packages)
func GetFilteredTickers(ctx context.Context, mcapMin *int64, inceptionMax *time.Time) ([]string, error) {
	// Handle nullable parameters
	var mcap int64
	if mcapMin != nil {
		mcap = *mcapMin
	}
	var inception time.Time
	if inceptionMax != nil {
		inception = *inceptionMax
	}

	tickers, err := genQ().GetFilteredTickers(ctx, generated.GetFilteredTickersParams{
		Column1:         types.PriceUpdateTypes,
		Column2:         mcap,
		Column3:         inception,
		LastPriceStatus: f.StringToNullString(types.StatusOK),
	})
	if err != nil {
		logf("[DB] Query error: %v\n", err)
		return nil, fmt.Errorf("failed to get filtered tickers: %w", err)
	}

	logf("[DB] GetFilteredTickers returned %d tickers\n", len(tickers))
	return tickers, nil
}

// GetTickersWithPrices returns tickers that have price data (limited to specified count)
func GetTickersWithPrices(ctx context.Context, limit int) ([]string, error) {
	tickers, err := genQ().GetTickersWithPrices(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get tickers with prices: %w", err)
	}
	return tickers, nil
}

// GetSymbolsWithStalePrices returns symbols with outdated price data
// Returns Symbol structs with only ticker and currency populated
func GetSymbolsWithStalePrices(ctx context.Context, limit int) ([]types.Symbol, error) {
	threshold := GetPriceThreshold()
	rows, err := genQ().GetSymbolsWithStalePrices(ctx, generated.GetSymbolsWithStalePricesParams{
		LastPriceUpdate: f.MaybeTimeToNullTime(&threshold),
		Column2:         types.PriceUpdateTypes,
		Limit:           int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get symbols with stale prices: %w", err)
	}

	var symbols []types.Symbol
	for _, row := range rows {
		symbols = append(symbols, types.Symbol{
			Ticker:   row.Ticker,
			Currency: f.NullStringToMaybeString(row.Currency),
		})
	}

	return symbols, nil
}

// GetPriceThreshold returns the threshold for stale prices (30 days ago)
// Changed from monthly to 30-day rolling window to reduce FMP queries
func GetPriceThreshold() time.Time {
	now := time.Now().UTC()
	return now.AddDate(0, 0, -30)
}

// CountStalePrices returns the count of stale prices
func CountStalePrices(ctx context.Context) (int, error) {
	threshold := GetPriceThreshold()
	count, err := genQ().CountStalePrices(ctx, generated.CountStalePricesParams{
		LastPriceUpdate: f.MaybeTimeToNullTime(&threshold),
		Column2:         types.PriceUpdateTypes,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count stale prices: %w", err)
	}
	return int(count), nil
}

// GetOldestPriceUpdate returns the oldest price update timestamp (only for actively trading symbols)
func GetOldestPriceUpdate(ctx context.Context) (*time.Time, error) {
	result, err := genQ().GetOldestPriceUpdate(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest price update: %w", err)
	}

	// Handle interface{} return from sqlc
	if result == nil {
		return nil, nil
	}
	if t, ok := result.(time.Time); ok {
		return &t, nil
	}
	return nil, nil
}

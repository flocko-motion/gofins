package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/f"
	"github.com/flocko-motion/gofins/pkg/types"
)

// GetSymbolsWithCIK returns all stock symbols that have a CIK value
func GetSymbolsWithCIK(ctx context.Context) ([]types.Symbol, error) {
	rows, err := genQ().GetSymbolsWithCIK(ctx, f.StringToNullString(types.TypeStock))
	if err != nil {
		return nil, fmt.Errorf("failed to get symbols with CIK: %w", err)
	}

	var symbols []types.Symbol
	for _, row := range rows {
		symbols = append(symbols, types.Symbol{
			Ticker:         row.Ticker,
			Exchange:       f.NullStringToMaybeString(row.Exchange),
			CIK:            f.NullStringToMaybeString(row.Cik),
			PrimaryListing: f.NullStringToMaybeString(row.PrimaryListing),
		})
	}

	return symbols, nil
}

// GetStockSymbolsWithoutCIK returns all stock symbols that don't have a CIK value
func GetStockSymbolsWithoutCIK(ctx context.Context) ([]types.Symbol, error) {
	rows, err := genQ().GetStockSymbolsWithoutCIK(ctx, f.StringToNullString(types.TypeStock))
	if err != nil {
		return nil, fmt.Errorf("failed to get stock symbols without CIK: %w", err)
	}

	var symbols []types.Symbol
	for _, row := range rows {
		symbols = append(symbols, types.Symbol{
			Ticker:         row.Ticker,
			Name:           f.NullStringToMaybeString(row.Name),
			OldestPrice:    f.NullTimeToMaybeTime(row.OldestPrice),
			PrimaryListing: f.NullStringToMaybeString(row.PrimaryListing),
		})
	}

	return symbols, nil
}

// GetStockSymbolsForNameDedupe returns all stock symbols with names for name-based deduplication
func GetStockSymbolsForNameDedupe(ctx context.Context) ([]types.Symbol, error) {
	rows, err := genQ().GetStockSymbolsForNameDedupe(ctx, f.StringToNullString(types.TypeStock))
	if err != nil {
		return nil, fmt.Errorf("failed to get stock symbols for name dedupe: %w", err)
	}

	var symbols []types.Symbol
	for _, row := range rows {
		symbols = append(symbols, types.Symbol{
			Ticker:         row.Ticker,
			Name:           f.NullStringToMaybeString(row.Name),
			OldestPrice:    f.NullTimeToMaybeTime(row.OldestPrice),
			PrimaryListing: f.NullStringToMaybeString(row.PrimaryListing),
		})
	}

	return symbols, nil
}

// ResetPrimaryListings sets all primary_listing fields to NULL
func ResetPrimaryListings(ctx context.Context) (int, error) {
	count, err := genQ().ResetPrimaryListings(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to reset primary listings: %w", err)
	}
	return int(count), nil
}

// UpdatePrimaryListingGroup updates primary_listing for a group of symbols
// primaryTicker gets primary_listing = ”, all others get primary_listing = primaryTicker
func UpdatePrimaryListingGroup(ctx context.Context, primaryTicker string, secondaryTickers []string) error {
	tx, err := Db().conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	q := genQ().WithTx(tx)

	// Set primary to empty string
	if err := q.SetPrimaryListing(ctx, primaryTicker); err != nil {
		return fmt.Errorf("failed to set primary listing: %w", err)
	}

	// Set all secondaries to point to primary
	if len(secondaryTickers) > 0 {
		if err := q.SetSecondaryListings(ctx, generated.SetSecondaryListingsParams{
			PrimaryListing: sql.NullString{String: primaryTicker, Valid: true},
			Column2:        secondaryTickers,
		}); err != nil {
			return fmt.Errorf("failed to set secondary listings: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

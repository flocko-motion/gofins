package db

import (
	"context"
	"fmt"
	"time"

	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/f"
)

// ErrorEntry represents an error logged during system operations
type ErrorEntry struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`    // e.g., "updater.symbols", "api.handler"
	ErrorType string    `json:"errorType"` // e.g., "network", "database", "validation"
	Message   string    `json:"message"`
	Details   *string   `json:"details,omitempty"`
}

// LogError logs an error to the database (implements types.ErrorLogger)
func (db *DB) LogError(source, level, message string, metadata map[string]interface{}) error {
	var details *string
	if metadata != nil {
		detailsStr := fmt.Sprintf("%v", metadata)
		details = &detailsStr
	}

	return genQ().InsertError(context.Background(), generated.InsertErrorParams{
		Source:    source,
		ErrorType: level,
		Message:   message,
		Details:   f.MaybeStringToNullString(details),
	})
}

// LogError is a package-level convenience function
func LogError(ctx context.Context, source, errorType, message string, details *string) error {
	fmt.Printf("[%s] %s: %s\n", source, errorType, message)
	return genQ().InsertError(ctx, generated.InsertErrorParams{
		Source:    source,
		ErrorType: errorType,
		Message:   message,
		Details:   f.MaybeStringToNullString(details),
	})
}

// GetRecentErrors retrieves the most recent errors
func GetRecentErrors(ctx context.Context, limit int) ([]ErrorEntry, error) {
	rows, err := genQ().GetRecentErrors(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get recent errors: %w", err)
	}

	var errors []ErrorEntry
	for _, row := range rows {
		timestamp := time.Time{}
		if ts := f.NullTimeToMaybeTime(row.Timestamp); ts != nil {
			timestamp = *ts
		}
		errors = append(errors, ErrorEntry{
			ID:        int(row.ID),
			Timestamp: timestamp,
			Source:    row.Source,
			ErrorType: row.ErrorType,
			Message:   row.Message,
			Details:   f.NullStringToMaybeString(row.Details),
		})
	}

	return errors, nil
}

// GetErrorsBySource retrieves errors for a specific source
func GetErrorsBySource(ctx context.Context, source string, limit int) ([]ErrorEntry, error) {
	rows, err := genQ().GetErrorsBySource(ctx, generated.GetErrorsBySourceParams{
		Source: source,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get errors by source: %w", err)
	}

	var errors []ErrorEntry
	for _, row := range rows {
		timestamp := time.Time{}
		if ts := f.NullTimeToMaybeTime(row.Timestamp); ts != nil {
			timestamp = *ts
		}
		errors = append(errors, ErrorEntry{
			ID:        int(row.ID),
			Timestamp: timestamp,
			Source:    row.Source,
			ErrorType: row.ErrorType,
			Message:   row.Message,
			Details:   f.NullStringToMaybeString(row.Details),
		})
	}

	return errors, nil
}

// CountErrorsSince counts errors since a given timestamp
func CountErrorsSince(ctx context.Context, since time.Time) (int, error) {
	count, err := genQ().CountErrorsSince(ctx, f.MaybeTimeToNullTime(&since))
	if err != nil {
		return 0, fmt.Errorf("failed to count errors since: %w", err)
	}
	return int(count), nil
}

// ClearOldErrors deletes errors older than the specified duration
func ClearOldErrors(ctx context.Context, olderThan time.Duration) (int, error) {
	cutoff := time.Now().Add(-olderThan)
	count, err := genQ().ClearOldErrors(ctx, f.MaybeTimeToNullTime(&cutoff))
	if err != nil {
		return 0, fmt.Errorf("failed to clear old errors: %w", err)
	}
	return int(count), nil
}

// GetErrorByID retrieves a specific error by its ID
func GetErrorByID(ctx context.Context, id int) (*ErrorEntry, error) {
	row, err := genQ().GetErrorByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get error by ID: %w", err)
	}

	timestamp := time.Time{}
	if ts := f.NullTimeToMaybeTime(row.Timestamp); ts != nil {
		timestamp = *ts
	}

	return &ErrorEntry{
		ID:        int(row.ID),
		Timestamp: timestamp,
		Source:    row.Source,
		ErrorType: row.ErrorType,
		Message:   row.Message,
		Details:   f.NullStringToMaybeString(row.Details),
	}, nil
}

// ClearAllErrors deletes all errors from the database
func ClearAllErrors(ctx context.Context) (int, error) {
	count, err := genQ().ClearAllErrors(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to clear all errors: %w", err)
	}
	return int(count), nil
}

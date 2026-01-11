package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/f"
)

type BatchUpdateLog struct {
	ID               int
	UpdaterName      string
	StartedAt        time.Time
	CompletedAt      *time.Time
	Status           string
	SymbolsProcessed int
	SymbolsUpdated   int
	ErrorMessage     *string
}

// StartBatchUpdate creates a new batch update log entry
func StartBatchUpdate(ctx context.Context, updaterName string) (int, error) {
	id, err := genQ().StartBatchUpdate(ctx, generated.StartBatchUpdateParams{
		UpdaterName: updaterName,
		StartedAt:   time.Now(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to start batch update log: %w", err)
	}
	return int(id), nil
}

// CompleteBatchUpdate marks a batch update as completed
func CompleteBatchUpdate(ctx context.Context, id int, symbolsProcessed, symbolsUpdated int) error {
	now := time.Now()
	err := genQ().CompleteBatchUpdate(ctx, generated.CompleteBatchUpdateParams{
		CompletedAt:      f.MaybeTimeToNullTime(&now),
		SymbolsProcessed: sql.NullInt32{Int32: int32(symbolsProcessed), Valid: true},
		SymbolsUpdated:   sql.NullInt32{Int32: int32(symbolsUpdated), Valid: true},
		ID:               int32(id),
	})
	if err != nil {
		return fmt.Errorf("failed to complete batch update log: %w", err)
	}
	return nil
}

// FailBatchUpdate marks a batch update as failed
func FailBatchUpdate(ctx context.Context, id int, errorMessage string) error {
	now := time.Now()
	err := genQ().FailBatchUpdate(ctx, generated.FailBatchUpdateParams{
		CompletedAt:  f.MaybeTimeToNullTime(&now),
		ErrorMessage: sql.NullString{String: errorMessage, Valid: true},
		ID:           int32(id),
	})
	if err != nil {
		return fmt.Errorf("failed to mark batch update as failed: %w", err)
	}
	return nil
}

// GetLastBatchUpdate returns the last batch update for a given updater
func GetLastBatchUpdate(ctx context.Context, updaterName string) (*BatchUpdateLog, error) {
	genLog, err := genQ().GetLastBatchUpdate(ctx, updaterName)
	if err != nil {
		return nil, err
	}

	var errorMessage *string
	if genLog.ErrorMessage.Valid {
		errorMessage = &genLog.ErrorMessage.String
	}

	return &BatchUpdateLog{
		ID:               int(genLog.ID),
		UpdaterName:      genLog.UpdaterName,
		StartedAt:        genLog.StartedAt,
		CompletedAt:      f.NullTimeToMaybeTime(genLog.CompletedAt),
		Status:           genLog.Status,
		SymbolsProcessed: int(genLog.SymbolsProcessed.Int32),
		SymbolsUpdated:   int(genLog.SymbolsUpdated.Int32),
		ErrorMessage:     errorMessage,
	}, nil
}

// DeleteBatchUpdate deletes all batch update logs for a given updater
func DeleteBatchUpdate(ctx context.Context, updaterName string) error {
	err := genQ().DeleteBatchUpdate(ctx, updaterName)
	if err != nil {
		return fmt.Errorf("failed to delete batch update logs: %w", err)
	}
	return nil
}

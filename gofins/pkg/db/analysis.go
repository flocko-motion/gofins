package db

import (
	"context"
	"database/sql"

	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/f"
	"github.com/flocko-motion/gofins/pkg/types"
	"github.com/google/uuid"
)

// CreateAnalysisPackage inserts a new analysis package with status='processing'
func CreateAnalysisPackage(ctx context.Context, pkg *types.AnalysisPackage) error {
	pkgUUID, err := uuid.Parse(pkg.ID)
	if err != nil {
		return err
	}

	return genQ().CreateAnalysisPackage(ctx, generated.CreateAnalysisPackageParams{
		ID:           pkgUUID,
		Name:         pkg.Name,
		CreatedAt:    pkg.CreatedAt,
		Interval:     pkg.Interval,
		TimeFrom:     pkg.TimeFrom,
		TimeTo:       pkg.TimeTo,
		HistBins:     int32(pkg.HistBins),
		HistMin:      pkg.HistMin,
		HistMax:      pkg.HistMax,
		McapMin:      f.MaybeInt64ToNullInt64(pkg.McapMin),
		InceptionMax: f.MaybeTimeToNullTime(pkg.InceptionMax),
		Status:       pkg.Status,
		UserID:       pkg.UserID,
	})
}

// UpdateAnalysisPackageStatus updates the status and symbol count of a package for a specific user
func UpdateAnalysisPackageStatus(ctx context.Context, userID uuid.UUID, packageID string, status string, symbolCount int) error {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return err
	}
	return genQ().UpdateAnalysisPackageStatus(ctx, generated.UpdateAnalysisPackageStatusParams{
		Status:      status,
		SymbolCount: sql.NullInt32{Int32: int32(symbolCount), Valid: true},
		ID:          pkgUUID,
		UserID:      userID,
	})
}

// SaveAnalysisResult saves a single analysis result (verifies package ownership)
func SaveAnalysisResult(ctx context.Context, userID uuid.UUID, packageID, ticker string, count int, mean, stddev, variance, min, max float64, histogramJSON []byte) error {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return err
	}

	// Verify package belongs to user
	exists, err := genQ().VerifyPackageOwnership(ctx, generated.VerifyPackageOwnershipParams{
		ID:     pkgUUID,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	return genQ().SaveAnalysisResult(ctx, generated.SaveAnalysisResultParams{
		PackageID: pkgUUID,
		Ticker:    ticker,
		Count:     int32(count),
		Mean:      mean,
		Stddev:    stddev,
		Variance:  variance,
		Min:       min,
		Max:       max,
		Histogram: histogramJSON,
	})
}

// GetAnalysisResults retrieves all results for a package (verifies package ownership)
func GetAnalysisResults(ctx context.Context, userID uuid.UUID, packageID string) ([]types.AnalysisResult, error) {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return nil, err
	}

	// Verify package belongs to user
	exists, err := genQ().VerifyPackageOwnership(ctx, generated.VerifyPackageOwnershipParams{
		ID:     pkgUUID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, sql.ErrNoRows
	}

	genResults, err := genQ().GetAnalysisResults(ctx, pkgUUID)
	if err != nil {
		return nil, err
	}

	results := make([]types.AnalysisResult, len(genResults))
	for i, r := range genResults {
		results[i] = types.AnalysisResult{
			PackageID:     r.PackageID.String(),
			Ticker:        r.Ticker,
			Count:         int(r.Count),
			Mean:          r.Mean,
			StdDev:        r.Stddev,
			Variance:      r.Variance,
			Min:           r.Min,
			Max:           r.Max,
			InceptionDate: f.NullTimeToMaybeTime(r.Inception),
		}
	}

	return results, nil
}

// GetAnalysisPackage retrieves a package by ID for a specific user
func GetAnalysisPackage(ctx context.Context, userID uuid.UUID, packageID string) (*types.AnalysisPackage, error) {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return nil, err
	}

	genPkg, err := genQ().GetAnalysisPackage(ctx, generated.GetAnalysisPackageParams{
		ID:     pkgUUID,
		UserID: userID,
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	symbolCount := 0
	if genPkg.SymbolCount.Valid {
		symbolCount = int(genPkg.SymbolCount.Int32)
	}

	return &types.AnalysisPackage{
		ID:           genPkg.ID.String(),
		Name:         genPkg.Name,
		CreatedAt:    genPkg.CreatedAt,
		Interval:     genPkg.Interval,
		TimeFrom:     genPkg.TimeFrom,
		TimeTo:       genPkg.TimeTo,
		HistBins:     int(genPkg.HistBins),
		HistMin:      genPkg.HistMin,
		HistMax:      genPkg.HistMax,
		McapMin:      f.NullInt64ToMaybeInt64(genPkg.McapMin),
		InceptionMax: f.NullTimeToMaybeTime(genPkg.InceptionMax),
		SymbolCount:  symbolCount,
		Status:       genPkg.Status,
		UserID:       genPkg.UserID,
	}, nil
}

// ListAnalysisPackages returns all analysis packages for a specific user
func ListAnalysisPackages(ctx context.Context, userID uuid.UUID) ([]types.AnalysisPackage, error) {
	genPkgs, err := genQ().ListAnalysisPackages(ctx, userID)
	if err != nil {
		return nil, err
	}

	packages := make([]types.AnalysisPackage, len(genPkgs))
	for i, genPkg := range genPkgs {
		symbolCount := 0
		if genPkg.SymbolCount.Valid {
			symbolCount = int(genPkg.SymbolCount.Int32)
		}

		packages[i] = types.AnalysisPackage{
			ID:           genPkg.ID.String(),
			Name:         genPkg.Name,
			CreatedAt:    genPkg.CreatedAt,
			Interval:     genPkg.Interval,
			TimeFrom:     genPkg.TimeFrom,
			TimeTo:       genPkg.TimeTo,
			HistBins:     int(genPkg.HistBins),
			HistMin:      genPkg.HistMin,
			HistMax:      genPkg.HistMax,
			McapMin:      f.NullInt64ToMaybeInt64(genPkg.McapMin),
			InceptionMax: f.NullTimeToMaybeTime(genPkg.InceptionMax),
			SymbolCount:  symbolCount,
			Status:       genPkg.Status,
			UserID:       genPkg.UserID,
		}
	}

	return packages, nil
}

// UpdateAnalysisPackageName updates the name of a package for a specific user
func UpdateAnalysisPackageName(ctx context.Context, userID uuid.UUID, packageID string, name string) error {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return err
	}
	return genQ().UpdateAnalysisPackageName(ctx, generated.UpdateAnalysisPackageNameParams{
		Name:   name,
		ID:     pkgUUID,
		UserID: userID,
	})
}

// DeleteAnalysisPackage deletes a package for a specific user (CASCADE will delete results)
func DeleteAnalysisPackage(ctx context.Context, userID uuid.UUID, packageID string) error {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return err
	}
	return genQ().DeleteAnalysisPackage(ctx, generated.DeleteAnalysisPackageParams{
		ID:     pkgUUID,
		UserID: userID,
	})
}

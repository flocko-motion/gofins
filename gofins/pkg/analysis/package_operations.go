package analysis

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flocko-motion/gofins/pkg/db"
	"github.com/flocko-motion/gofins/pkg/types"
	"github.com/google/uuid"
)

// GetPackage retrieves a package by ID
func GetPackage(ctx context.Context, userID uuid.UUID, packageID string) (*types.AnalysisPackage, error) {
	return db.GetAnalysisPackage(ctx, userID, packageID)
}

// ListPackages returns all analysis packages
func ListPackages(ctx context.Context, userID uuid.UUID) ([]types.AnalysisPackage, error) {
	return db.ListAnalysisPackages(ctx, userID)
}

// UpdatePackageName updates the name of an analysis package
func UpdatePackageName(ctx context.Context, userID uuid.UUID, packageID string, name string) (*types.AnalysisPackage, error) {
	// Check if package exists
	pkg, err := db.GetAnalysisPackage(ctx, userID, packageID)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, fmt.Errorf("package not found")
	}

	// Update in database
	if err := db.UpdateAnalysisPackageName(ctx, userID, packageID, name); err != nil {
		return nil, err
	}

	// Return updated package
	pkg.Name = name
	return pkg, nil
}

// DeletePackage deletes an analysis package and its associated files
func DeletePackage(ctx context.Context, userID uuid.UUID, packageID string) error {
	// First, check if package exists
	pkg, err := db.GetAnalysisPackage(ctx, userID, packageID)
	if err != nil {
		return fmt.Errorf("failed to get package: %w", err)
	}
	if pkg == nil {
		return fmt.Errorf("package not found")
	}

	// Delete from database first
	if err := db.DeleteAnalysisPackage(ctx, userID, packageID); err != nil {
		return fmt.Errorf("failed to delete package from database: %w", err)
	}

	// Clean up PNG files if they exist
	plotsDir := PathPlots(packageID)
	if err := cleanupPackageFiles(packageID, plotsDir); err != nil {
		// Log the error but don't fail the deletion since DB is already cleaned
		errMsg := fmt.Sprintf("Failed to cleanup files for package %s: %v", packageID, err)
		_ = db.LogError("analysis.package", "filesystem", "Failed to cleanup package files", &errMsg)
	}

	return nil
}

// cleanupPackageFiles removes all files associated with a package
func cleanupPackageFiles(packageID, plotsDir string) error {
	// Check if plots directory exists
	if _, err := os.Stat(plotsDir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, nothing to clean up
	}

	// Find all PNG files in the plots directory
	pattern := filepath.Join(plotsDir, "*.png")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to find PNG files: %w", err)
	}

	// Delete each PNG file
	deletedCount := 0
	for _, file := range matches {
		if err := os.Remove(file); err != nil {
			errMsg := fmt.Sprintf("Failed to delete file %s: %v", file, err)
			_ = db.LogError("analysis.package", "filesystem", "Failed to delete analysis file", &errMsg)
		} else {
			deletedCount++
		}
	}

	// Try to remove the entire directory if it's empty
	if err := os.Remove(plotsDir); err != nil {
		// Directory not empty or other error - that's fine, we cleaned the files
		fmt.Printf("[ANALYSIS] Note: Could not remove directory %s (may not be empty)\n", plotsDir)
	} else {
		fmt.Printf("[ANALYSIS] Removed plots directory for package %s\n", packageID)
	}

	fmt.Printf("[ANALYSIS] Cleaned up %d PNG files for package %s\n", deletedCount, packageID)
	return nil
}

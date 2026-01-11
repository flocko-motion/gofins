package reset

import (
	"fmt"

	"github.com/flocko-motion/gofins/pkg/db"
	"github.com/spf13/cobra"
)

var symbolCmd = &cobra.Command{
	Use:   "symbol [ticker]",
	Short: "Reset price and profile timestamps for a single symbol",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticker := args[0]

		fmt.Printf("=== Resetting Symbol: %s ===\n", ticker)

		err := db.ResetSymbol(cmd.Context(), ticker)
		if err != nil {
			return fmt.Errorf("failed to reset symbol: %w", err)
		}

		fmt.Printf("✓ Reset timestamps for %s\n", ticker)
		fmt.Println("Updater will reload price and profile data from FMP on next run")

		return nil
	},
}

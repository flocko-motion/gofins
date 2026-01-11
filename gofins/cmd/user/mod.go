package user

import (
	"fmt"

	"github.com/flocko-motion/gofins/pkg/db"
	"github.com/spf13/cobra"
)

var (
	adminFlag *bool
)

var modCmd = &cobra.Command{
	Use:   "mod [name-or-uuid]",
	Short: "Modify user properties",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]

		// Check if admin flag was provided
		if adminFlag == nil {
			return fmt.Errorf("no flags provided, use --admin to set admin status")
		}

		// Update user admin status
		user, err := db.UpdateUserAdmin(cmd.Context(), nameOrID, *adminFlag)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		adminStr := "no"
		if user.IsAdmin {
			adminStr = "yes"
		}

		fmt.Printf("✓ User updated:\n")
		fmt.Printf("  Name:    %s\n", user.Name)
		fmt.Printf("  ID:      %s\n", user.ID)
		fmt.Printf("  Admin:   %s\n", adminStr)
		fmt.Printf("  Created: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

func init() {
	adminFlag = modCmd.Flags().Bool("admin", false, "Set admin status (true/false)")
	UserCmd.AddCommand(modCmd)
}

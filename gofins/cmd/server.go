package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flocko-motion/gofins/pkg/api"
	"github.com/flocko-motion/gofins/pkg/db"
	"github.com/flocko-motion/gofins/pkg/f"
	"github.com/flocko-motion/gofins/pkg/updater"
	"github.com/spf13/cobra"
)

var noUpdates bool
var devUser string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run FINS in server mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Setup context for graceful shutdown
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// init DB - initialize singleton with retry logic
		fmt.Println("=== Initializing ===")
		var database *db.DB
		maxRetries := 30
		retryDelay := time.Second
		
		for i := 0; i < maxRetries; i++ {
			database = db.Db()
			if database != nil {
				break
			}
			
			if i == 0 {
				fmt.Println("⏳ Database not ready, retrying...")
			}
			
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				// Exponential backoff, max 10 seconds
				retryDelay = time.Duration(float64(retryDelay) * 1.5)
				if retryDelay > 10*time.Second {
					retryDelay = 10 * time.Second
				}
				// Reset dbOnce to allow retry
				db.ResetConnection()
			}
		}
		
		if database == nil {
			return fmt.Errorf("failed to connect to database after %d retries", maxRetries)
		}
		fmt.Println("✓ Database connected")

		// Start REST API server
		apiServer := api.NewServer(db.Db(), 8080, devUser)
		go apiServer.Start(ctx)
		if devUser != "" {
			fmt.Printf("✓ REST API server listening on :8080 (DEV MODE - all requests as user '%s')\n", devUser)
		} else {
			fmt.Println("✓ REST API server listening on :8080")
		}

		// Start updaters only if --no-updates is not set
		fmt.Println("\n=== Starting Services ===")
		if !noUpdates {
			go updater.RunAllUpdaters(ctx)
		} else {
			fmt.Println("⚠️  Updates disabled - working with existing data only")
		}

		fmt.Println("✓ Server running (Ctrl+C to stop)")
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
		time.Sleep(time.Second)

		// Graceful shutdown
		if err := db.PrepareForShutdown(); err != nil {
			_ = db.LogError("server.shutdown", "database", "Failed to close database connection", f.Ptr(err.Error()))
		}

		fmt.Println("✓ Server stopped")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolVar(&noUpdates, "no-updates", false,
		"Disable all data updates and work with existing data only")
	serverCmd.Flags().StringVar(&devUser, "user", "",
		"Development mode - override user for all requests (e.g., --user=alice)")
}

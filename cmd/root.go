package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "otg-sport-api",
	Short: "OTG Sport API - Sports data sync and API service",
	Long: `OTG Sport API is a service that syncs sports data from GoalServe
and exposes it through a REST API.

Available commands:
  serve  - Start the REST API server
  sync   - Run the data sync scheduler
  apikey - Manage API keys`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

package cmd

import (
	"github.com/dusanbre/otg-sports-api/cmd/commands"
	"github.com/spf13/cobra"
)

var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Manage API keys",
	Long:  `Create, list, and revoke API keys for accessing the REST API.`,
}

func init() {
	rootCmd.AddCommand(apikeyCmd)

	// Add subcommands
	apikeyCmd.AddCommand(commands.ApiKeyCreateCmd)
	apikeyCmd.AddCommand(commands.ApiKeyListCmd)
	apikeyCmd.AddCommand(commands.ApiKeyRevokeCmd)
}

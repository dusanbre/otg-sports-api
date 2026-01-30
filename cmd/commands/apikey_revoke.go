package commands

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var ApiKeyRevokeCmd = &cobra.Command{
	Use:   "revoke [id]",
	Short: "Revoke an API key",
	Long: `Revoke an API key by its ID. The key will be deactivated and can no longer be used.

Example:
  otg-sport-api apikey revoke 5`,
	Args: cobra.ExactArgs(1),
	Run:  runApiKeyRevoke,
}

func runApiKeyRevoke(cmd *cobra.Command, args []string) {
	// Parse ID
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Fatalf("Invalid ID: %s", args[0])
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get database instance
	db, err := database.GetInstance()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Revoke the key
	if err := db.RevokeApiKey(id); err != nil {
		log.Fatalf("Failed to revoke API key: %v", err)
	}

	fmt.Printf("\n✓ API key with ID %d has been revoked.\n\n", id)
}

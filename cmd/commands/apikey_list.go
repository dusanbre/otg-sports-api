package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var ApiKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	Long:  `List all API keys with their details (key prefix, name, sports, status).`,
	Run:   runApiKeyList,
}

func runApiKeyList(cmd *cobra.Command, args []string) {
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

	// Fetch all API keys
	keys, err := db.GetAllApiKeys()
	if err != nil {
		log.Fatalf("Failed to fetch API keys: %v", err)
	}

	if len(keys) == 0 {
		fmt.Println("No API keys found.")
		return
	}

	// Print header
	fmt.Println()
	fmt.Printf("%-6s %-15s %-25s %-25s %-10s %-10s\n", "ID", "PREFIX", "NAME", "SPORTS", "RATE", "STATUS")
	fmt.Println(strings.Repeat("-", 95))

	// Print each key
	for _, k := range keys {
		status := "active"
		if !k.IsActive {
			status = "revoked"
		}

		sports := strings.Join(k.Sports, ",")
		if len(sports) > 23 {
			sports = sports[:20] + "..."
		}

		name := k.Name
		if len(name) > 23 {
			name = name[:20] + "..."
		}

		fmt.Printf("%-6d %-15s %-25s %-25s %-10d %-10s\n",
			k.ID, k.KeyPrefix, name, sports, k.RateLimit, status)
	}
	fmt.Println()
}

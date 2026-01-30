package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/dusanbre/otg-sports-api/internal/api/auth"
	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	keyName   string
	keySports string
	rateLimit int
)

var ApiKeyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Long: `Create a new API key with specified name and sport access.
Examples:
  otg-sport-api apikey create --name "Mobile App" --sports soccer,basketball
  otg-sport-api apikey create --name "Admin Key" --sports "*"
  otg-sport-api apikey create --name "Soccer Only" --sports soccer --rate-limit 200`,
	Run: runApiKeyCreate,
}

func init() {
	ApiKeyCreateCmd.Flags().StringVarP(&keyName, "name", "n", "", "Name/description for the API key (required)")
	ApiKeyCreateCmd.Flags().StringVarP(&keySports, "sports", "s", "", "Comma-separated sports (soccer,basketball) or * for all (required)")
	ApiKeyCreateCmd.Flags().IntVarP(&rateLimit, "rate-limit", "r", 100, "Rate limit (requests per minute)")
	ApiKeyCreateCmd.MarkFlagRequired("name")
	ApiKeyCreateCmd.MarkFlagRequired("sports")
}

func runApiKeyCreate(cmd *cobra.Command, args []string) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Parse sports
	var sports []string
	if keySports == "*" {
		sports = []string{"*"}
	} else {
		sports = strings.Split(keySports, ",")
		for i, s := range sports {
			sports[i] = strings.TrimSpace(s)
		}
	}

	// Validate sports
	validSports := map[string]bool{"soccer": true, "basketball": true, "*": true}
	for _, s := range sports {
		if !validSports[s] {
			log.Fatalf("Invalid sport: %s. Valid options: soccer, basketball, *", s)
		}
	}

	// Get database instance
	db, err := database.GetInstance()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Generate API key
	plainKey, keyHash, keyPrefix, err := auth.GenerateAPIKey()
	if err != nil {
		log.Fatalf("Failed to generate API key: %v", err)
	}

	// Insert into database
	_, err = db.CreateApiKey(keyHash, keyPrefix, keyName, sports, rateLimit)
	if err != nil {
		log.Fatalf("Failed to create API key: %v", err)
	}

	// Print success message
	fmt.Println()
	fmt.Println("✓ API key created successfully")
	fmt.Println()
	fmt.Printf("  Name:       %s\n", keyName)
	fmt.Printf("  Sports:     %s\n", strings.Join(sports, ", "))
	fmt.Printf("  Rate Limit: %d req/min\n", rateLimit)
	fmt.Printf("  Key:        %s\n", plainKey)
	fmt.Println()
	fmt.Println("  ⚠️  Save this key now - it won't be shown again!")
	fmt.Println()
}

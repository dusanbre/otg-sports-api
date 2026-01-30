package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dusanbre/otg-sports-api/internal/api"
	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	apiPort string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the REST API server",
	Long: `Start the REST API server that exposes sports data.

The server provides endpoints for:
  - Soccer matches (GET /api/v1/soccer/matches)
  - Basketball matches (GET /api/v1/basketball/matches)
  - Live matches for each sport
  - League listings

Authentication is required via API key:
  - Header: Authorization: Bearer <api_key>
  - Header: X-API-Key: <api_key>`,
	Run: runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&apiPort, "port", "p", "8080", "Port to run the API server on")
}

func runServe(cmd *cobra.Command, args []string) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Override port from env if set
	if envPort := os.Getenv("API_PORT"); envPort != "" {
		apiPort = envPort
	}

	// Get database instance
	db, err := database.GetInstance()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to database!")

	// Create and start server
	server := api.NewServer(db, apiPort)

	// Start server in goroutine
	go func() {
		fmt.Printf("Starting API server on port %s...\n", apiPort)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give 30 seconds for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped")
}

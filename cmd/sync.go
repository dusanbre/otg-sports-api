package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/dusanbre/otg-sports-api/internal/services"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Run the data sync scheduler",
	Long: `Start the data sync scheduler that periodically fetches sports data
from GoalServe and stores it in the database.

The scheduler runs every minute and syncs:
  - Soccer matches (today and next 7 days)
  - Basketball matches (today and next 7 days)`,
	Run: runSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) {
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

	fmt.Println("Successfully connected to database!")

	// Create sync services
	soccerSyncService := services.NewSoccerSyncService(db)
	basketballSyncService := services.NewBasketballSyncService(db)

	// Create scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	// Schedule soccer sync job
	soccerJob, err := scheduler.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(func() {
			log.Println("Running scheduled soccer match sync...")
			if err := soccerSyncService.SyncMatches(); err != nil {
				log.Printf("Error syncing soccer matches: %v", err)
			}
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create soccer job: %v", err)
	}
	fmt.Printf("Scheduled soccer job with ID: %s - runs every 1 minute\n", soccerJob.ID())

	// Schedule basketball sync job
	basketballJob, err := scheduler.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(func() {
			log.Println("Running scheduled basketball match sync...")
			if err := basketballSyncService.SyncMatches(); err != nil {
				log.Printf("Error syncing basketball matches: %v", err)
			}
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create basketball job: %v", err)
	}
	fmt.Printf("Scheduled basketball job with ID: %s - runs every 1 minute\n", basketballJob.ID())

	// Run initial sync
	log.Println("Running initial soccer match sync...")
	if err := soccerSyncService.SyncMatches(); err != nil {
		log.Printf("Error in initial soccer sync: %v", err)
	}

	log.Println("Running initial basketball match sync...")
	if err := basketballSyncService.SyncMatches(); err != nil {
		log.Printf("Error in initial basketball sync: %v", err)
	}

	// Start scheduler
	scheduler.Start()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down scheduler...")
	if err := scheduler.Shutdown(); err != nil {
		log.Printf("Error shutting down scheduler: %v", err)
	}

	log.Println("Scheduler stopped")
}

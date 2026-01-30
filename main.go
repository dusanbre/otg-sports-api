package main

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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database instance (singleton pattern)
	db, err := database.GetInstance()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to database!")

	// Create soccer sync service
	soccerSyncService := services.NewSoccerSyncService(db)

	// Create a new scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	// Schedule the soccer match sync job to run every 1 minute
	job, err := scheduler.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(func() {
			log.Println("Running scheduled soccer match sync...")
			if err := soccerSyncService.SyncMatches(); err != nil {
				log.Printf("Error syncing soccer matches: %v", err)
			}
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create job: %v", err)
	}

	fmt.Printf("Scheduled job with ID: %s - runs every 1 minute\n", job.ID())

	// Run the sync immediately on startup
	log.Println("Running initial soccer match sync...")
	if err := soccerSyncService.SyncMatches(); err != nil {
		log.Printf("Error in initial soccer sync: %v", err)
	}

	// Start the scheduler
	scheduler.Start()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down scheduler...")
	if err := scheduler.Shutdown(); err != nil {
		log.Printf("Error shutting down scheduler: %v", err)
	}

	log.Println("Scheduler stopped")
}

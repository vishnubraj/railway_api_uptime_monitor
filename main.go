package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"railway-api-uptime-monitor/internal/config"
	"railway-api-uptime-monitor/internal/database"
	"railway-api-uptime-monitor/internal/monitor"
	"railway-api-uptime-monitor/internal/server"
	"railway-api-uptime-monitor/internal/webhook"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect()

	// Initialize webhook notifier
	notifier := webhook.NewNotifier(cfg)

	// Initialize monitor
	apiMonitor := monitor.New(db, notifier, cfg)

	// Set up cron job for monitoring
	c := cron.New()
	_, err = c.AddFunc(cfg.CheckInterval, func() {
		log.Println("Running scheduled API health checks...")
		apiMonitor.CheckAllAPIs()
	})
	if err != nil {
		log.Fatalf("Failed to set up cron job: %v", err)
	}
	c.Start()

	// Initialize and start web server
	srv := server.New(db, cfg)

	// Graceful shutdown
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Stop cron jobs
	c.Stop()

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Simple health checker for Railway cron job
// This can be run as a separate service or cron job on Railway

type APIConfig struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	Method         string `json:"method"`
	ExpectedStatus int    `json:"expected_status"`
	Timeout        int    `json:"timeout"`
}

type APIsConfig struct {
	APIs []APIConfig `json:"apis"`
}

type HealthCheck struct {
	APIName      string        `bson:"api_name"`
	URL          string        `bson:"url"`
	Status       string        `bson:"status"`
	StatusCode   int           `bson:"status_code"`
	ResponseTime time.Duration `bson:"response_time"`
	Timestamp    time.Time     `bson:"timestamp"`
	ErrorMessage string        `bson:"error_message,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cron.go <check|alert>")
	}

	command := os.Args[1]

	switch command {
	case "check":
		runHealthChecks()
	case "alert":
		checkForAlerts()
	default:
		log.Fatal("Unknown command. Use 'check' or 'alert'")
	}
}

func runHealthChecks() {
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("DATABASE_NAME", "uptime_monitor")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Load API configuration
	apisConfig, err := loadAPIs()
	if err != nil {
		log.Fatalf("Failed to load APIs: %v", err)
	}

	// Perform health checks
	collection := client.Database(dbName).Collection("health_checks")

	for _, apiConfig := range apisConfig.APIs {
		start := time.Now()
		status, statusCode, err := performHealthCheck(apiConfig)
		responseTime := time.Since(start)

		healthCheck := HealthCheck{
			APIName:      apiConfig.Name,
			URL:          apiConfig.URL,
			Status:       status,
			StatusCode:   statusCode,
			ResponseTime: responseTime,
			Timestamp:    time.Now(),
		}

		if err != nil {
			healthCheck.ErrorMessage = err.Error()
		}

		_, insertErr := collection.InsertOne(context.Background(), healthCheck)
		if insertErr != nil {
			log.Printf("Error inserting health check: %v", insertErr)
		}

		log.Printf("Checked %s: %s (%d) - %v", apiConfig.Name, status, statusCode, responseTime)
	}
}

func checkForAlerts() {
	// This function would check for APIs that need alerts
	// and send notifications via webhooks
	log.Println("Checking for alerts...")
	// Implementation would go here
}

func loadAPIs() (*APIsConfig, error) {
	file, err := os.ReadFile("config/apis.json")
	if err != nil {
		return nil, err
	}

	var config APIsConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func performHealthCheck(apiConfig APIConfig) (string, int, error) {
	client := &http.Client{
		Timeout: time.Duration(apiConfig.Timeout) * time.Second,
	}

	req, err := http.NewRequest(apiConfig.Method, apiConfig.URL, nil)
	if err != nil {
		return "down", 0, err
	}

	req.Header.Set("User-Agent", "Railway-API-Uptime-Monitor-Cron/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "down", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == apiConfig.ExpectedStatus {
		return "up", resp.StatusCode, nil
	}

	return "down", resp.StatusCode, fmt.Errorf("unexpected status code: %d, expected: %d", resp.StatusCode, apiConfig.ExpectedStatus)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

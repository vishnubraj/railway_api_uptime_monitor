package monitor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"railway-api-uptime-monitor/internal/config"
	"railway-api-uptime-monitor/internal/database"
	"railway-api-uptime-monitor/internal/models"
	"railway-api-uptime-monitor/internal/webhook"

	"go.mongodb.org/mongo-driver/bson"
)

type Monitor struct {
	db       *database.Database
	notifier *webhook.Notifier
	config   *config.Config
	client   *http.Client
}

func New(db *database.Database, notifier *webhook.Notifier, cfg *config.Config) *Monitor {
	return &Monitor{
		db:       db,
		notifier: notifier,
		config:   cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		},
	}
}

func (m *Monitor) CheckAllAPIs() {
	apisConfig, err := config.LoadAPIs()
	if err != nil {
		log.Printf("Error loading API config: %v", err)
		return
	}

	for _, apiConfig := range apisConfig.APIs {
		go m.checkAPI(apiConfig)
	}
}

func (m *Monitor) checkAPI(apiConfig config.APIConfig) {
	start := time.Now()

	status, statusCode, err := m.performHealthCheck(apiConfig)
	responseTime := time.Since(start)

	healthCheck := models.HealthCheck{
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.GetCollection("health_checks")
	_, insertErr := collection.InsertOne(ctx, healthCheck)
	if insertErr != nil {
		log.Printf("Error inserting health check: %v", insertErr)
	}

	m.updateAPIStatus(apiConfig, status, statusCode, responseTime, err)

	log.Printf("Checked %s: %s (%d) - %v", apiConfig.Name, status, statusCode, responseTime)
}

func (m *Monitor) performHealthCheck(apiConfig config.APIConfig) (string, int, error) {
	var req *http.Request
	var err error

	if apiConfig.Method == "POST" {
		req, err = http.NewRequest("POST", apiConfig.URL, bytes.NewBuffer([]byte{}))
	} else {
		req, err = http.NewRequest("GET", apiConfig.URL, nil)
	}

	if err != nil {
		return "down", 0, err
	}

	req.Header.Set("User-Agent", "Railway-API-Uptime-Monitor/1.0")

	resp, err := m.client.Do(req)
	if err != nil {
		return "down", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == apiConfig.ExpectedStatus {
		return "up", resp.StatusCode, nil
	}

	return "down", resp.StatusCode, fmt.Errorf("unexpected status code: %d, expected: %d", resp.StatusCode, apiConfig.ExpectedStatus)
}

func (m *Monitor) updateAPIStatus(apiConfig config.APIConfig, status string, statusCode int, responseTime time.Duration, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.GetCollection("api_status")

	var existingStatus models.APIStatus
	filter := bson.M{"name": apiConfig.Name}
	findErr := collection.FindOne(ctx, filter).Decode(&existingStatus)

	now := time.Now()

	if findErr != nil {
		newStatus := models.APIStatus{
			Name:          apiConfig.Name,
			URL:           apiConfig.URL,
			Method:        apiConfig.Method,
			Status:        status,
			StatusCode:    statusCode,
			ResponseTime:  responseTime,
			LastChecked:   now,
			DowntimeCount: 0,
			UptimePercent: 100.0,
		}

		if status == "up" {
			newStatus.LastUp = now
		} else {
			newStatus.LastDown = now
			newStatus.DowntimeCount = 1
			newStatus.UptimePercent = 0.0
		}

		if err != nil {
			newStatus.ErrorMessage = err.Error()
		}

		_, insertErr := collection.InsertOne(ctx, newStatus)
		if insertErr != nil {
			log.Printf("Error inserting API status: %v", insertErr)
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"status":        status,
			"status_code":   statusCode,
			"response_time": responseTime,
			"last_checked":  now,
		},
	}

	if status == "up" {
		update["$set"].(bson.M)["last_up"] = now
		update["$set"].(bson.M)["error_message"] = ""

		if existingStatus.Status == "down" {
			update["$set"].(bson.M)["downtime_count"] = 0
			m.sendAlert(apiConfig.Name, "up", "API is back online")
		}
	} else {
		update["$set"].(bson.M)["last_down"] = now
		if err != nil {
			update["$set"].(bson.M)["error_message"] = err.Error()
		}

		newDowntimeCount := existingStatus.DowntimeCount + 1
		update["$set"].(bson.M)["downtime_count"] = newDowntimeCount

		if newDowntimeCount >= m.config.DowntimeThreshold {
			message := fmt.Sprintf("API has been down for %d consecutive checks", newDowntimeCount)
			m.sendAlert(apiConfig.Name, "down", message)
		}
	}

	uptimePercent := m.calculateUptimePercent(apiConfig.Name)
	update["$set"].(bson.M)["uptime_percent"] = uptimePercent

	_, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		log.Printf("Error updating API status: %v", updateErr)
	}
}

func (m *Monitor) calculateUptimePercent(apiName string) float64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.GetCollection("health_checks")

	since := time.Now().Add(-24 * time.Hour)
	filter := bson.M{
		"api_name":  apiName,
		"timestamp": bson.M{"$gte": since},
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil || total == 0 {
		return 100.0
	}

	upFilter := bson.M{
		"api_name":  apiName,
		"timestamp": bson.M{"$gte": since},
		"status":    "up",
	}

	upCount, err := collection.CountDocuments(ctx, upFilter)
	if err != nil {
		return 0.0
	}

	return float64(upCount) / float64(total) * 100.0
}

func (m *Monitor) sendAlert(apiName, alertType, message string) {
	alert := models.Alert{
		APIName:   apiName,
		Type:      alertType,
		Message:   message,
		Timestamp: time.Now(),
		Resolved:  alertType == "up",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.GetCollection("alerts")
	_, err := collection.InsertOne(ctx, alert)
	if err != nil {
		log.Printf("Error storing alert: %v", err)
	}

	go m.notifier.SendAlert(apiName, alertType, message)
}

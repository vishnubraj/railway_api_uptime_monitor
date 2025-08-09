package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"railway-api-uptime-monitor/internal/database"
	"railway-api-uptime-monitor/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	db *database.Database
}

func New(db *database.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "railway-api-uptime-monitor",
	})
}

func (h *Handler) Dashboard(c *gin.Context) {
	// Get all API statuses for dashboard
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.db.GetCollection("api_status")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{
			"error": "Failed to load API statuses",
		})
		return
	}
	defer cursor.Close(ctx)

	var apiStatuses []models.APIStatus
	if err := cursor.All(ctx, &apiStatuses); err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{
			"error": "Failed to decode API statuses",
		})
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"apis":      apiStatuses,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

func (h *Handler) GetAllStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.db.GetCollection("api_status")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var apiStatuses []models.APIStatus
	if err := cursor.All(ctx, &apiStatuses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apis":      apiStatuses,
		"timestamp": time.Now(),
	})
}

func (h *Handler) GetAPIStatus(c *gin.Context) {
	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := h.db.GetCollection("api_status")
	var apiStatus models.APIStatus
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&apiStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
		return
	}

	c.JSON(http.StatusOK, apiStatus)
}

func (h *Handler) GetAPILogs(c *gin.Context) {
	name := c.Param("name")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.db.GetCollection("health_checks")
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"timestamp": -1})

	cursor, err := collection.Find(ctx, bson.M{"api_name": name}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var healthChecks []models.HealthCheck
	if err := cursor.All(ctx, &healthChecks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"api_name": name,
		"logs":     healthChecks,
		"count":    len(healthChecks),
	})
}

func (h *Handler) GetAlerts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	unresolvedOnly := c.Query("unresolved") == "true"

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.db.GetCollection("alerts")

	filter := bson.M{}
	if unresolvedOnly {
		filter["resolved"] = false
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"timestamp": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var alerts []models.Alert
	if err := cursor.All(ctx, &alerts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"count":  len(alerts),
	})
}

func (h *Handler) GetStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get total APIs
	apiCollection := h.db.GetCollection("api_status")
	totalAPIs, _ := apiCollection.CountDocuments(ctx, bson.M{})

	// Get APIs that are up
	upAPIs, _ := apiCollection.CountDocuments(ctx, bson.M{"status": "up"})

	// Get total checks in last 24 hours
	checksCollection := h.db.GetCollection("health_checks")
	since := time.Now().Add(-24 * time.Hour)
	totalChecks, _ := checksCollection.CountDocuments(ctx, bson.M{
		"timestamp": bson.M{"$gte": since},
	})

	// Get unresolved alerts
	alertsCollection := h.db.GetCollection("alerts")
	unresolvedAlerts, _ := alertsCollection.CountDocuments(ctx, bson.M{"resolved": false})

	c.JSON(http.StatusOK, gin.H{
		"total_apis":        totalAPIs,
		"apis_up":           upAPIs,
		"apis_down":         totalAPIs - upAPIs,
		"total_checks_24h":  totalChecks,
		"unresolved_alerts": unresolvedAlerts,
		"uptime_percentage": func() float64 {
			if totalAPIs == 0 {
				return 100.0
			}
			return float64(upAPIs) / float64(totalAPIs) * 100.0
		}(),
		"timestamp": time.Now(),
	})
}

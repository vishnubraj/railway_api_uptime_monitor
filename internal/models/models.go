package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APIStatus struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	URL           string             `bson:"url" json:"url"`
	Method        string             `bson:"method" json:"method"`
	Status        string             `bson:"status" json:"status"` // "up", "down", "unknown"
	StatusCode    int                `bson:"status_code" json:"status_code"`
	ResponseTime  time.Duration      `bson:"response_time" json:"response_time"`
	LastChecked   time.Time          `bson:"last_checked" json:"last_checked"`
	LastUp        time.Time          `bson:"last_up" json:"last_up"`
	LastDown      time.Time          `bson:"last_down" json:"last_down"`
	DowntimeCount int                `bson:"downtime_count" json:"downtime_count"`
	UptimePercent float64            `bson:"uptime_percent" json:"uptime_percent"`
	ErrorMessage  string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
}

type HealthCheck struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIName      string             `bson:"api_name" json:"api_name"`
	URL          string             `bson:"url" json:"url"`
	Status       string             `bson:"status" json:"status"`
	StatusCode   int                `bson:"status_code" json:"status_code"`
	ResponseTime time.Duration      `bson:"response_time" json:"response_time"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
	ErrorMessage string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
}

type Alert struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIName   string             `bson:"api_name" json:"api_name"`
	Type      string             `bson:"type" json:"type"` // "down", "up", "timeout"
	Message   string             `bson:"message" json:"message"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Resolved  bool               `bson:"resolved" json:"resolved"`
}

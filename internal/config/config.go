package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port              string
	MongoURI          string
	DatabaseName      string
	CheckInterval     string
	TimeoutSeconds    int
	MaxRetries        int
	SlackWebhookURL   string
	DiscordWebhookURL string
	EnableSlack       bool
	EnableDiscord     bool
	DowntimeThreshold int
}

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

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		MongoURI:          getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName:      getEnv("DATABASE_NAME", "uptime_monitor"),
		CheckInterval:     getEnv("CHECK_INTERVAL", "*/5 * * * *"),
		TimeoutSeconds:    getEnvAsInt("TIMEOUT_SECONDS", 30),
		MaxRetries:        getEnvAsInt("MAX_RETRIES", 3),
		SlackWebhookURL:   getEnv("SLACK_WEBHOOK_URL", ""),
		DiscordWebhookURL: getEnv("DISCORD_WEBHOOK_URL", ""),
		EnableSlack:       getEnvAsBool("ENABLE_SLACK", false),
		EnableDiscord:     getEnvAsBool("ENABLE_DISCORD", false),
		DowntimeThreshold: getEnvAsInt("DOWNTIME_THRESHOLD", 3),
	}
}

func LoadAPIs() (*APIsConfig, error) {
	file, err := os.ReadFile("config/apis.json")
	if err != nil {
		// Create default config if file doesn't exist
		defaultConfig := &APIsConfig{
			APIs: []APIConfig{
				{
					Name:           "Example API",
					URL:            "https://httpbin.org/status/200",
					Method:         "GET",
					ExpectedStatus: 200,
					Timeout:        30,
				},
			},
		}

		// Create config directory and file
		os.MkdirAll("config", 0755)
		defaultJSON, _ := json.MarshalIndent(defaultConfig, "", "  ")
		os.WriteFile("config/apis.json", defaultJSON, 0644)

		return defaultConfig, nil
	}

	var config APIsConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Invalid boolean value for %s: %s, using default: %t", key, value, defaultValue)
	}
	return defaultValue
}

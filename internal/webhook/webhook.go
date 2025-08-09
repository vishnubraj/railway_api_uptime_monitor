package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"railway-api-uptime-monitor/internal/config"
)

type Notifier struct {
	config *config.Config
	client *http.Client
}

type SlackPayload struct {
	Text        string            `json:"text"`
	Username    string            `json:"username"`
	IconEmoji   string            `json:"icon_emoji"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Color     string       `json:"color"`
	Title     string       `json:"title"`
	Text      string       `json:"text"`
	Fields    []SlackField `json:"fields,omitempty"`
	Timestamp int64        `json:"ts"`
}

type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type DiscordPayload struct {
	Content   string         `json:"content"`
	Username  string         `json:"username"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Timestamp   string              `json:"timestamp"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func NewNotifier(cfg *config.Config) *Notifier {
	return &Notifier{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (n *Notifier) SendAlert(apiName, alertType, message string) {
	if n.config.EnableSlack && n.config.SlackWebhookURL != "" {
		go n.sendSlackAlert(apiName, alertType, message)
	}

	if n.config.EnableDiscord && n.config.DiscordWebhookURL != "" {
		go n.sendDiscordAlert(apiName, alertType, message)
	}
}

func (n *Notifier) sendSlackAlert(apiName, alertType, message string) {
	color := "#ff0000" // Red for down
	emoji := ":x:"
	if alertType == "up" {
		color = "#00ff00" // Green for up
		emoji = ":white_check_mark:"
	}

	payload := SlackPayload{
		Text:      fmt.Sprintf("%s API Alert", apiName),
		Username:  "Railway Uptime Monitor",
		IconEmoji: ":robot_face:",
		Attachments: []SlackAttachment{
			{
				Color: color,
				Title: fmt.Sprintf("%s %s API Status Change", emoji, apiName),
				Text:  message,
				Fields: []SlackField{
					{
						Title: "API Name",
						Value: apiName,
						Short: true,
					},
					{
						Title: "Status",
						Value: alertType,
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().Format("2006-01-02 15:04:05 UTC"),
						Short: false,
					},
				},
				Timestamp: time.Now().Unix(),
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling Slack payload: %v", err)
		return
	}

	resp, err := n.client.Post(n.config.SlackWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending Slack notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Slack webhook returned status: %d", resp.StatusCode)
	} else {
		log.Printf("Slack notification sent for %s: %s", apiName, alertType)
	}
}

func (n *Notifier) sendDiscordAlert(apiName, alertType, message string) {
	color := 16711680 // Red for down
	if alertType == "up" {
		color = 65280 // Green for up
	}

	payload := DiscordPayload{
		Username: "Railway Uptime Monitor",
		Embeds: []DiscordEmbed{
			{
				Title:       fmt.Sprintf("%s API Status Change", apiName),
				Description: message,
				Color:       color,
				Fields: []DiscordEmbedField{
					{
						Name:   "API Name",
						Value:  apiName,
						Inline: true,
					},
					{
						Name:   "Status",
						Value:  alertType,
						Inline: true,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling Discord payload: %v", err)
		return
	}

	resp, err := n.client.Post(n.config.DiscordWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending Discord notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		log.Printf("Discord webhook returned status: %d", resp.StatusCode)
	} else {
		log.Printf("Discord notification sent for %s: %s", apiName, alertType)
	}
}

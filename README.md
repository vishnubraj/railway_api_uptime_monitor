# Railway API Uptime Monitor

A robust API uptime monitoring service built with Go, MongoDB, and designed for Railway deployment with cron jobs and webhook integrations.

## Features

- 🔄 Periodic API health checks with configurable intervals
- 📊 MongoDB storage for status logs and historical data
- 📈 Web dashboard with real-time status monitoring
- ⚠️ Downtime alerts with customizable thresholds
- 🔗 Slack/Discord webhook integration
- ⚙️ Fully configurable via environment variables
- 🚀 Railway-ready with cron jobs and database services

## Tech Stack

- **Backend**: Go (Gin framework)
- **Database**: MongoDB
- **Scheduling**: Cron jobs
- **Deployment**: Railway
- **Notifications**: Webhook integration (Slack/Discord)

## Environment Variables

Create a `.env` file with the following variables:

```env
# Server Configuration
PORT=8080
GIN_MODE=release

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=uptime_monitor

# Monitoring Configuration
CHECK_INTERVAL=*/5 * * * *  # Every 5 minutes
TIMEOUT_SECONDS=30
MAX_RETRIES=3

# Webhook Configuration
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK
ENABLE_SLACK=true
ENABLE_DISCORD=false

# Alert Configuration
DOWNTIME_THRESHOLD=3  # Number of consecutive failures before alert
```

## API Endpoints Configuration

Configure the APIs to monitor in `config/apis.json`:

```json
{
  "apis": [
    {
      "name": "Example API",
      "url": "https://api.example.com/health",
      "method": "GET",
      "expected_status": 200,
      "timeout": 30
    }
  ]
}
```

## Railway Deployment

1. Connect your repository to Railway
2. Add the required environment variables in Railway dashboard
3. Enable MongoDB addon in Railway
4. Deploy!

The app will automatically:
- Set up cron jobs for monitoring
- Create database collections
- Start the web dashboard

## API Endpoints

- `GET /` - Dashboard
- `GET /api/status` - Current status of all monitored APIs
- `GET /api/logs/:name` - Historical logs for a specific API
- `GET /api/health` - Service health check

## Local Development

```bash
# Install dependencies
go mod tidy

# Set up environment
cp .env.example .env
# Edit .env with your configuration

# Run the application
go run main.go
```

## License

MIT

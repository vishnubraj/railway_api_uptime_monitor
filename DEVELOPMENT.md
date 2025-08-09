# API Uptime Monitor - Development Guide

## Project Structure

```
railway-api-uptime-monitor/
├── main.go                     # Main application entry point
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── Dockerfile                  # Container configuration
├── railway.toml               # Railway deployment configuration
├── deploy.sh                  # Deployment script
├── .env.example              # Environment variables template
├── README.md                 # Project documentation
├── cmd/
│   └── cron/
│       └── main.go          # Standalone cron job binary
├── config/
│   └── apis.json           # API endpoints configuration
├── internal/
│   ├── config/
│   │   └── config.go       # Configuration management
│   ├── database/
│   │   └── database.go     # MongoDB connection
│   ├── handlers/
│   │   └── handlers.go     # HTTP request handlers
│   ├── models/
│   │   └── models.go       # Data models
│   ├── monitor/
│   │   └── monitor.go      # API monitoring logic
│   ├── server/
│   │   └── server.go       # HTTP server setup
│   └── webhook/
│       └── webhook.go      # Notification webhooks
└── web/
    └── templates/
        └── dashboard.html  # Web dashboard template
```

## Development Setup

### Prerequisites

- Go 1.21 or later
- MongoDB (local or cloud)
- Git

### Local Development

1. **Clone and setup**
   ```bash
   cd railway-api-uptime-monitor
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure APIs**
   Edit `config/apis.json` to add your APIs to monitor:
   ```json
   {
     "apis": [
       {
         "name": "My API",
         "url": "https://api.example.com/health",
         "method": "GET",
         "expected_status": 200,
         "timeout": 30
       }
     ]
   }
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

5. **Access the dashboard**
   Open http://localhost:8080

### Building

```bash
# Build main application
go build -o uptime-monitor .

# Build cron job binary
go build -o cron-job ./cmd/cron
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DATABASE_NAME` | Database name | `uptime_monitor` |
| `CHECK_INTERVAL` | Cron schedule for checks | `*/5 * * * *` |
| `TIMEOUT_SECONDS` | HTTP request timeout | `30` |
| `MAX_RETRIES` | Max retry attempts | `3` |
| `SLACK_WEBHOOK_URL` | Slack webhook URL | - |
| `DISCORD_WEBHOOK_URL` | Discord webhook URL | - |
| `ENABLE_SLACK` | Enable Slack notifications | `false` |
| `ENABLE_DISCORD` | Enable Discord notifications | `false` |
| `DOWNTIME_THRESHOLD` | Failures before alert | `3` |

### API Configuration

Edit `config/apis.json`:

```json
{
  "apis": [
    {
      "name": "API Name",
      "url": "https://api.example.com/endpoint",
      "method": "GET|POST",
      "expected_status": 200,
      "timeout": 30
    }
  ]
}
```

## Deployment

### Railway

1. **Install Railway CLI**
   ```bash
   npm install -g @railway/cli
   ```

2. **Deploy using script**
   ```bash
   ./deploy.sh
   ```

3. **Manual deployment**
   ```bash
   railway login
   railway init
   railway up
   ```

4. **Add MongoDB service**
   - Go to Railway dashboard
   - Add MongoDB service
   - Copy connection string to `MONGODB_URI`

### Docker

```bash
# Build image
docker build -t uptime-monitor .

# Run container
docker run -d \
  -p 8080:8080 \
  -e MONGODB_URI=mongodb://localhost:27017 \
  uptime-monitor
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Dashboard UI |
| `/api/health` | GET | Service health check |
| `/api/status` | GET | All API statuses |
| `/api/status/:name` | GET | Specific API status |
| `/api/logs/:name` | GET | API health check logs |
| `/api/alerts` | GET | Recent alerts |
| `/api/stats` | GET | System statistics |

## Database Schema

### Collections

#### `api_status`
```javascript
{
  _id: ObjectId,
  name: String,
  url: String,
  method: String,
  status: String,        // "up", "down", "unknown"
  status_code: Number,
  response_time: Number, // in milliseconds
  last_checked: Date,
  last_up: Date,
  last_down: Date,
  downtime_count: Number,
  uptime_percent: Number,
  error_message: String
}
```

#### `health_checks`
```javascript
{
  _id: ObjectId,
  api_name: String,
  url: String,
  status: String,
  status_code: Number,
  response_time: Number,
  timestamp: Date,
  error_message: String
}
```

#### `alerts`
```javascript
{
  _id: ObjectId,
  api_name: String,
  type: String,     // "down", "up", "timeout"
  message: String,
  timestamp: Date,
  resolved: Boolean
}
```

## Monitoring Features

- **Health Checks**: Periodic API monitoring with configurable intervals
- **Status Tracking**: Real-time status updates and historical data
- **Uptime Calculation**: 24-hour rolling uptime percentage
- **Alert System**: Configurable downtime threshold alerts
- **Webhook Notifications**: Slack and Discord integration
- **Web Dashboard**: Real-time status visualization
- **REST API**: Programmatic access to monitoring data

## Webhook Integration

### Slack

1. Create a Slack app and webhook
2. Set `SLACK_WEBHOOK_URL` and `ENABLE_SLACK=true`

### Discord

1. Create a Discord webhook in your server
2. Set `DISCORD_WEBHOOK_URL` and `ENABLE_DISCORD=true`

## Troubleshooting

### Common Issues

1. **Build fails**: Ensure Go modules are enabled with `GO111MODULE=on`
2. **MongoDB connection**: Check `MONGODB_URI` format and network access
3. **API not found**: Verify `config/apis.json` exists and is valid JSON
4. **Webhooks not working**: Check webhook URLs and enable flags

### Logs

Check application logs for detailed error information:
```bash
# Railway
railway logs

# Local
go run main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

# Railway API Uptime Monitor Template

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/railway-api-uptime-monitor)

## 🚀 One-Click Deployment

This template provides a complete API uptime monitoring solution that can be deployed to Railway with just one click.

### What You Get

- ✅ **Go-based uptime monitor** with MongoDB storage
- ✅ **Real-time web dashboard** with beautiful UI
- ✅ **Automated health checks** every 5 minutes
- ✅ **Slack/Discord notifications** for downtime alerts
- ✅ **REST API** for programmatic access
- ✅ **Historical data tracking** with uptime percentages
- ✅ **Configurable alert thresholds**

## 📋 Pre-configured Settings

The template comes with these default configurations:

| Setting | Default Value | Description |
|---------|---------------|-------------|
| Check Interval | Every 5 minutes | How often APIs are checked |
| Timeout | 30 seconds | Request timeout for health checks |
| Alert Threshold | 3 failures | Consecutive failures before alert |
| Port | 8080 | Application port |
| Database | MongoDB | Included MongoDB service |

## 🎯 Quick Start (3 Steps)

### 1. Deploy Template
Click the "Deploy on Railway" button above or use the deployment script:
```bash
curl -sL https://raw.githubusercontent.com/username/railway-api-uptime-monitor/main/railway-deploy-template.sh | bash
```

### 2. Configure APIs
After deployment, edit `config/apis.json` in your repository:
```json
{
  "apis": [
    {
      "name": "My Important API",
      "url": "https://api.myservice.com/health",
      "method": "GET",
      "expected_status": 200,
      "timeout": 30
    }
  ]
}
```

### 3. Set Up Notifications (Optional)
Add webhook URLs in Railway dashboard variables:
```bash
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
ENABLE_SLACK=true
```

## 🔧 Environment Variables

The template automatically sets up these variables:

### Required (Auto-configured)
- `MONGODB_URI` - MongoDB connection string
- `DATABASE_NAME` - Database name
- `PORT` - Application port

### Optional (Customizable)
- `CHECK_INTERVAL` - Cron schedule (default: `*/5 * * * *`)
- `TIMEOUT_SECONDS` - Request timeout (default: `30`)
- `DOWNTIME_THRESHOLD` - Alert threshold (default: `3`)
- `SLACK_WEBHOOK_URL` - Slack notifications
- `DISCORD_WEBHOOK_URL` - Discord notifications
- `ENABLE_SLACK` - Enable Slack alerts
- `ENABLE_DISCORD` - Enable Discord alerts

## 📊 Features Included

### Dashboard
- Real-time API status monitoring
- Uptime percentage calculations
- Response time tracking
- Error message display
- Auto-refresh every 30 seconds

### Alerts
- Configurable downtime thresholds
- Recovery notifications
- Webhook integration (Slack/Discord)
- Alert history tracking

### API Endpoints
- `GET /` - Web dashboard
- `GET /api/status` - All API statuses
- `GET /api/status/:name` - Specific API status
- `GET /api/logs/:name` - Health check history
- `GET /api/alerts` - Alert history
- `GET /api/stats` - System statistics

## 🛠️ Customization

### Adding More APIs
1. Edit `config/apis.json` in your repository
2. Commit changes
3. Railway will automatically redeploy

### Changing Check Frequency
Update the `CHECK_INTERVAL` environment variable:
```bash
railway variables set CHECK_INTERVAL="*/1 * * * *"  # Every minute
```

### Custom Notifications
Modify `internal/webhook/webhook.go` for custom message formats.

### UI Customization
Edit `web/templates/dashboard.html` for custom styling.

## 📁 Template Structure

```
railway-api-uptime-monitor/
├── main.go                    # Application entry point
├── railway.toml              # Railway configuration
├── railway-template.json     # Template definition
├── config/apis.json          # API endpoints to monitor
├── internal/                 # Application code
│   ├── config/              # Configuration management
│   ├── database/            # MongoDB connection
│   ├── handlers/            # HTTP handlers
│   ├── models/              # Data models
│   ├── monitor/             # Monitoring logic
│   ├── server/              # HTTP server
│   └── webhook/             # Notification system
└── web/templates/           # Dashboard UI
```

## 🔍 Monitoring Examples

### Basic Health Check
```json
{
  "name": "API Health",
  "url": "https://api.example.com/health",
  "method": "GET",
  "expected_status": 200,
  "timeout": 30
}
```

### Custom Endpoint
```json
{
  "name": "User Service",
  "url": "https://api.example.com/users/1",
  "method": "GET",
  "expected_status": 200,
  "timeout": 15
}
```

### POST Endpoint
```json
{
  "name": "Auth Service",
  "url": "https://api.example.com/auth/verify",
  "method": "POST",
  "expected_status": 200,
  "timeout": 45
}
```

## 🚨 Troubleshooting

### Common Issues

**MongoDB Connection Failed**
- Check that MongoDB service is running in Railway
- Verify `MONGODB_URI` environment variable

**APIs Not Being Checked**
- Verify `config/apis.json` is valid JSON
- Check application logs: `railway logs`

**Webhooks Not Working**
- Verify webhook URLs are correct
- Check `ENABLE_SLACK` or `ENABLE_DISCORD` is set to `true`

### Getting Help

1. Check Railway logs: `railway logs`
2. View environment variables: `railway variables`
3. Check service status in Railway dashboard

## 📈 Advanced Configuration

### Custom Cron Schedule
```bash
# Every minute
railway variables set CHECK_INTERVAL="* * * * *"

# Every hour
railway variables set CHECK_INTERVAL="0 * * * *"

# Every day at 9 AM
railway variables set CHECK_INTERVAL="0 9 * * *"
```

### Multiple Environments
Deploy separate instances for staging/production:
```bash
railway environment create staging
railway environment use staging
railway up
```

## 🎉 What's Next?

After successful deployment:

1. **Monitor your APIs** - Add your endpoints to `config/apis.json`
2. **Set up alerts** - Configure Slack/Discord webhooks
3. **Customize thresholds** - Adjust alert sensitivity
4. **Scale monitoring** - Add more APIs as needed
5. **Integrate with CI/CD** - Automate deployments

## 📄 License

MIT License - See LICENSE file for details.

---

**Deploy your API monitoring solution in under 2 minutes! 🚀**

# 🚀 Railway API Uptime Monitor - Quick Start

## What's Included

✅ **Complete Go Application** with MongoDB integration  
✅ **Web Dashboard** with real-time status monitoring  
✅ **Cron Job Monitoring** with configurable intervals  
✅ **Webhook Notifications** (Slack & Discord)  
✅ **Railway Deployment Ready** with all configurations  
✅ **Docker Support** for containerized deployment  
✅ **REST API** for programmatic access  

## 🎯 Quick Start (3 Steps)

### 1. Configure Your APIs
Edit `config/apis.json`:
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

### 2. Set Environment Variables
Copy and edit the environment file:
```bash
cp .env.example .env
# Edit .env with your MongoDB and webhook URLs
```

### 3. Deploy to Railway
```bash
./deploy.sh
```

That's it! Your uptime monitor is now running! 🎉

## 📊 Features At a Glance

| Feature | Description | Configuration |
|---------|-------------|---------------|
| **API Monitoring** | Periodic health checks | `CHECK_INTERVAL` in .env |
| **Downtime Alerts** | Configurable failure threshold | `DOWNTIME_THRESHOLD` in .env |
| **Slack Notifications** | Real-time alerts to Slack | `SLACK_WEBHOOK_URL` + `ENABLE_SLACK=true` |
| **Discord Notifications** | Real-time alerts to Discord | `DISCORD_WEBHOOK_URL` + `ENABLE_DISCORD=true` |
| **Web Dashboard** | Beautiful real-time UI | Access at your Railway URL |
| **Historical Data** | 24-hour uptime tracking | Stored in MongoDB |
| **REST API** | Programmatic access | `/api/*` endpoints |

## 🔧 Local Development

```bash
# Setup
make setup

# Run locally
make run

# Build
make build

# View dashboard
open http://localhost:8080
```

## 🌐 Railway Deployment

### Automatic Deployment
```bash
./deploy.sh
```

### Manual Deployment
1. Install Railway CLI: `npm install -g @railway/cli`
2. Login: `railway login`
3. Initialize: `railway init`
4. Add MongoDB service in Railway dashboard
5. Set environment variables in Railway dashboard
6. Deploy: `railway up`

## 📡 Webhook Setup

### Slack Integration
1. Create a Slack app at https://api.slack.com/apps
2. Add incoming webhook
3. Copy webhook URL to `SLACK_WEBHOOK_URL`
4. Set `ENABLE_SLACK=true`

### Discord Integration
1. Go to your Discord server settings
2. Integrations → Webhooks → New Webhook
3. Copy webhook URL to `DISCORD_WEBHOOK_URL`
4. Set `ENABLE_DISCORD=true`

## 📈 API Endpoints

- `GET /` - Web dashboard
- `GET /api/status` - All API statuses
- `GET /api/status/:name` - Specific API status
- `GET /api/logs/:name` - Historical logs
- `GET /api/alerts` - Recent alerts
- `GET /api/stats` - System statistics

## 🎛️ Configuration Files

| File | Purpose |
|------|---------|
| `.env` | Environment variables |
| `config/apis.json` | APIs to monitor |
| `railway.toml` | Railway deployment config |
| `Dockerfile` | Container configuration |

## 🔍 Monitoring Your APIs

The system automatically:
1. **Checks APIs** based on your cron schedule
2. **Stores results** in MongoDB
3. **Calculates uptime** over 24-hour periods
4. **Sends alerts** when downtime threshold is reached
5. **Tracks recovery** and sends "back online" notifications

## 🛠️ Customization

### Adding New APIs
Edit `config/apis.json` and restart the application.

### Changing Check Frequency
Update `CHECK_INTERVAL` in your environment (cron format).

### Webhook Message Format
Modify `internal/webhook/webhook.go` for custom message formatting.

### Dashboard Styling
Edit `web/templates/dashboard.html` for custom UI.

## 📚 Documentation

- **README.md** - Overview and setup
- **DEVELOPMENT.md** - Detailed development guide
- **This file** - Quick start guide

## 🆘 Need Help?

1. Check the logs: `railway logs` (Railway) or console output (local)
2. Verify environment variables are set correctly
3. Ensure MongoDB is accessible
4. Check that `config/apis.json` is valid JSON

## 📄 License

MIT License - See LICENSE file for details.

---

**Made with ❤️ for Railway deployment**

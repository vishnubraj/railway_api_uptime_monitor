# Railway Deployment Template

This repository includes comprehensive Railway deployment templates and configurations for easy one-click deployment.

## 🚀 Quick Deploy Options

### Option 1: One-Click Deploy Template
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/railway-api-uptime-monitor)

### Option 2: Automated Script Deploy
```bash
curl -sL https://raw.githubusercontent.com/username/railway-api-uptime-monitor/main/railway-deploy-template.sh | bash
```

### Option 3: Manual Railway CLI Deploy
```bash
# Clone repository
git clone https://github.com/username/railway-api-uptime-monitor.git
cd railway-api-uptime-monitor

# Deploy using Railway CLI
./railway-deploy-template.sh
```

## 📋 Template Configurations

The template includes multiple configuration formats:

### `railway.toml` - Basic Railway Configuration
```toml
[build]
builder = "NIXPACKS"

[deploy]
startCommand = "./uptime-monitor"
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 10

[variables]
PORT = "8080"
GIN_MODE = "release"
```

### `railway.json` - Advanced Configuration
```json
{
  "build": {
    "builder": "NIXPACKS",
    "buildCommand": "go build -o uptime-monitor ."
  },
  "deploy": {
    "startCommand": "./uptime-monitor",
    "restartPolicyType": "ON_FAILURE"
  }
}
```

### `railway-template.yml` - Multi-Service Template
Includes both the uptime monitor and MongoDB service with automatic connection configuration.

## 🛠️ Template Features

### Automatically Configured
- ✅ **Go Build Environment** - Optimized Nixpacks configuration
- ✅ **MongoDB Service** - Pre-configured database with persistence
- ✅ **Environment Variables** - All required settings pre-configured
- ✅ **Health Checks** - Automatic service monitoring
- ✅ **Auto-restart** - Failure recovery with retry policy

### Pre-configured Settings
| Setting | Value | Description |
|---------|-------|-------------|
| Port | 8080 | Application port |
| Check Interval | Every 5 minutes | API monitoring frequency |
| Timeout | 30 seconds | Request timeout |
| Alert Threshold | 3 failures | Before notifications |
| Restart Policy | On Failure | Auto-recovery enabled |

### Optional Features (Easy to Enable)
- 🔔 **Slack Notifications** - Just add webhook URL
- 🔔 **Discord Notifications** - Just add webhook URL  
- 📊 **Custom Dashboards** - Modify templates
- 🔧 **API Configuration** - Edit JSON file

## 📝 Post-Deployment Setup

### 1. Configure APIs to Monitor
Edit `config/apis.json`:
```json
{
  "apis": [
    {
      "name": "Production API",
      "url": "https://api.myservice.com/health",
      "method": "GET",
      "expected_status": 200,
      "timeout": 30
    },
    {
      "name": "User Service",
      "url": "https://users.myservice.com/ping",
      "method": "GET", 
      "expected_status": 200,
      "timeout": 15
    }
  ]
}
```

### 2. Set Up Notifications
In Railway dashboard → Variables:
```bash
# Slack
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
ENABLE_SLACK=true

# Discord  
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
ENABLE_DISCORD=true
```

### 3. Customize Monitoring
```bash
# Check every minute instead of 5 minutes
CHECK_INTERVAL=* * * * *

# More sensitive alerts (alert after 1 failure)
DOWNTIME_THRESHOLD=1

# Longer timeout for slow APIs
TIMEOUT_SECONDS=60
```

## 🔧 Template Customization

### For Template Publishers
1. Fork this repository
2. Update GitHub URLs in template configurations
3. Customize default settings as needed
4. Submit to Railway template gallery

### For Users
1. Deploy template
2. Modify `config/apis.json` for your APIs
3. Set environment variables for notifications
4. Customize UI in `web/templates/` if desired

## 📊 What's Included

### Services
- **Uptime Monitor** - Main Go application
- **MongoDB** - Database with persistent storage

### Features
- **Real-time Dashboard** - Web UI with status visualization
- **API Monitoring** - Configurable health checks
- **Alert System** - Threshold-based notifications
- **Historical Data** - 24-hour uptime tracking
- **REST API** - Programmatic access to data
- **Webhook Integration** - Slack/Discord notifications

### Monitoring Capabilities
- HTTP/HTTPS endpoint monitoring
- Custom status code validation
- Response time tracking
- Failure count and recovery detection
- Uptime percentage calculation
- Alert management and history

## 🚨 Template Deployment Troubleshooting

### Common Issues

**Build Fails**
```bash
# Check Go module configuration
railway logs --tail 100
```

**MongoDB Connection Issues**
```bash
# Verify MongoDB service is running
railway ps
# Check connection string
railway variables get MONGODB_URI
```

**Environment Variables Not Set**
```bash
# View all variables
railway variables
# Set missing variables
railway variables set VARIABLE_NAME=value
```

### Useful Commands
```bash
# View deployment logs
railway logs

# Check service status  
railway ps

# Open dashboard
railway open

# Redeploy
railway up --detach
```

## 🔗 Template Links

- **Railway Template Page**: https://railway.app/template/railway-api-uptime-monitor
- **Deploy Button**: `[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/railway-api-uptime-monitor)`
- **Source Repository**: https://github.com/username/railway-api-uptime-monitor
- **Documentation**: See README.md and DEVELOPMENT.md

## 📈 Template Versions

### v1.0 - Initial Release
- Basic uptime monitoring
- MongoDB integration
- Web dashboard

### v1.1 - Enhanced Features  
- Slack/Discord notifications
- Improved error handling
- Better UI/UX

### v1.2 - Template Optimization
- Optimized Railway configurations
- Better default settings
- Enhanced documentation

## 📄 License

MIT License - Template and application are free to use and modify.

---

**Deploy your API monitoring solution in under 2 minutes! 🚀**

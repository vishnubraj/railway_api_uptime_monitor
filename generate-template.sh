#!/bin/bash

# Railway Template Generator
# This script creates all necessary files for Railway template deployment

echo "🏗️  Generating Railway Deployment Template"
echo "=========================================="

# Create template directory structure
mkdir -p .railway-template/{configs,scripts,docs}

# Generate template configuration
cat > .railway-template/template.json << 'EOF'
{
  "name": "API Uptime Monitor",
  "description": "A comprehensive API uptime monitoring service with real-time dashboard, alerts, and webhook notifications. Built with Go, MongoDB, and designed for Railway deployment.",
  "tags": ["monitoring", "uptime", "api", "alerts", "dashboard", "go", "mongodb"],
  "categories": ["Monitoring", "DevOps", "SaaS"],
  "githubUrl": "https://github.com/username/railway-api-uptime-monitor",
  "demoUrl": "https://uptime-monitor-demo.railway.app",
  "image": "https://raw.githubusercontent.com/username/railway-api-uptime-monitor/main/.railway-template/preview.png",
  "services": [
    {
      "name": "uptime-monitor",
      "type": "web",
      "plan": "hobby",
      "source": {
        "type": "repo",
        "repo": "username/railway-api-uptime-monitor",
        "ref": "main"
      },
      "variables": {
        "PORT": {
          "description": "Application port",
          "default": "8080"
        },
        "GIN_MODE": {
          "description": "Gin framework mode",
          "default": "release"
        },
        "CHECK_INTERVAL": {
          "description": "Cron schedule for API checks",
          "default": "*/5 * * * *"
        },
        "TIMEOUT_SECONDS": {
          "description": "HTTP request timeout in seconds",
          "default": "30"
        },
        "DOWNTIME_THRESHOLD": {
          "description": "Number of consecutive failures before alert",
          "default": "3"
        },
        "SLACK_WEBHOOK_URL": {
          "description": "Slack webhook URL for notifications (optional)",
          "optional": true
        },
        "DISCORD_WEBHOOK_URL": {
          "description": "Discord webhook URL for notifications (optional)",
          "optional": true
        },
        "ENABLE_SLACK": {
          "description": "Enable Slack notifications",
          "default": "false"
        },
        "ENABLE_DISCORD": {
          "description": "Enable Discord notifications",
          "default": "false"
        }
      }
    },
    {
      "name": "mongodb",
      "type": "database",
      "plan": "hobby",
      "source": {
        "type": "image",
        "image": "mongo:7.0"
      },
      "variables": {
        "MONGO_INITDB_ROOT_USERNAME": {
          "description": "MongoDB root username",
          "default": "admin"
        },
        "MONGO_INITDB_ROOT_PASSWORD": {
          "description": "MongoDB root password",
          "generate": true
        },
        "MONGO_INITDB_DATABASE": {
          "description": "Initial database name",
          "default": "uptime_monitor"
        }
      },
      "volumes": [
        {
          "name": "mongodb_data",
          "mountPath": "/data/db"
        }
      ]
    }
  ],
  "connections": [
    {
      "from": "uptime-monitor",
      "to": "mongodb",
      "variable": "MONGODB_URI",
      "template": "mongodb://admin:${MONGO_INITDB_ROOT_PASSWORD}@${RAILWAY_PRIVATE_DOMAIN}:27017/uptime_monitor?authSource=admin"
    }
  ]
}
EOF

# Generate setup instructions
cat > .railway-template/SETUP.md << 'EOF'
# Railway Template Setup Guide

## Automatic Deployment

1. Click "Deploy Template" on Railway
2. Configure your APIs in the deployed repository
3. Set up notifications (optional)

## Post-Deployment Configuration

### 1. Configure APIs to Monitor

Edit `config/apis.json` in your deployed repository:

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

### 2. Set Up Notifications (Optional)

In Railway dashboard, set these environment variables:

**For Slack:**
- `SLACK_WEBHOOK_URL`: Your Slack webhook URL
- `ENABLE_SLACK`: Set to `true`

**For Discord:**
- `DISCORD_WEBHOOK_URL`: Your Discord webhook URL
- `ENABLE_DISCORD`: Set to `true`

### 3. Access Your Dashboard

Your uptime monitor will be available at your Railway app URL.

## Features Included

- ✅ Real-time API monitoring
- ✅ Web dashboard with status visualization
- ✅ Configurable alert thresholds
- ✅ Slack/Discord notifications
- ✅ Historical data tracking
- ✅ REST API for integration
- ✅ MongoDB data persistence

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CHECK_INTERVAL` | Cron schedule for checks | `*/5 * * * *` |
| `TIMEOUT_SECONDS` | Request timeout | `30` |
| `DOWNTIME_THRESHOLD` | Failures before alert | `3` |
| `SLACK_WEBHOOK_URL` | Slack notifications | Optional |
| `DISCORD_WEBHOOK_URL` | Discord notifications | Optional |

## Customization

After deployment, you can:
- Add/remove APIs in `config/apis.json`
- Adjust check frequency via `CHECK_INTERVAL`
- Customize alert thresholds
- Modify the dashboard UI
- Add custom notification channels

## Support

If you need help:
1. Check Railway logs
2. Review the documentation
3. Open an issue on GitHub
EOF

# Generate Railway button HTML
cat > .railway-template/deploy-button.html << 'EOF'
<a href="https://railway.app/template/railway-api-uptime-monitor?referralCode=template">
  <img src="https://railway.app/button.svg" alt="Deploy on Railway" />
</a>
EOF

# Generate Railway button Markdown
cat > .railway-template/deploy-button.md << 'EOF'
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/railway-api-uptime-monitor?referralCode=template)
EOF

# Create preview image placeholder
cat > .railway-template/preview-image.md << 'EOF'
# Preview Image

Create a screenshot of your dashboard and save it as `preview.png` in this directory.

Recommended dimensions: 1200x630 pixels

The image should show:
- The main dashboard with API status cards
- Some example APIs being monitored
- The clean, professional UI

This image will be used as the template preview in Railway's template gallery.
EOF

echo "✅ Railway template files generated in .railway-template/"
echo ""
echo "📝 Next steps:"
echo "1. Create a preview image (1200x630px) and save as .railway-template/preview.png"
echo "2. Update GitHub URLs in template.json with your repository"
echo "3. Submit your template to Railway"
echo ""
echo "🔗 Template files created:"
echo "   - .railway-template/template.json (main template config)"
echo "   - .railway-template/SETUP.md (user setup guide)"
echo "   - .railway-template/deploy-button.html (deploy button)"
echo "   - .railway-template/deploy-button.md (deploy button markdown)"

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

#!/bin/bash

# Railway Deployment Script
# This script helps deploy the API Uptime Monitor to Railway

echo "🚀 Railway API Uptime Monitor Deployment"
echo "========================================"

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "❌ Railway CLI is not installed. Please install it first:"
    echo "   npm install -g @railway/cli"
    exit 1
fi

# Login to Railway (if not already logged in)
echo "🔐 Checking Railway authentication..."
if ! railway whoami &> /dev/null; then
    echo "Please login to Railway:"
    railway login
fi

# Initialize project if needed
if [ ! -f ".railway/project.json" ]; then
    echo "📦 Initializing Railway project..."
    railway init
fi

# Set up environment variables
echo "⚙️  Setting up environment variables..."

# Check if .env file exists
if [ -f ".env" ]; then
    echo "📄 Found .env file. Please manually set these variables in Railway dashboard:"
    echo "   Go to: https://railway.app/dashboard"
    echo "   Select your project → Variables tab"
    echo ""
    echo "Required variables from .env:"
    grep -v '^#\|^$' .env
else
    echo "📝 Creating .env from template..."
    cp .env.example .env
    echo "✏️  Please edit .env with your configuration, then run this script again"
    exit 0
fi

# Add MongoDB service
echo "🗄️  Setting up MongoDB..."
echo "Please add MongoDB service in Railway dashboard:"
echo "   1. Go to your project dashboard"
echo "   2. Click 'New Service' → 'Database' → 'Add MongoDB'"
echo "   3. Once added, copy the connection string to MONGODB_URI in your variables"

# Deploy
echo "🚀 Deploying to Railway..."
railway up

echo "✅ Deployment complete!"
echo ""
echo "📋 Next steps:"
echo "   1. Set up your environment variables in Railway dashboard"
echo "   2. Add MongoDB service if not already done"
echo "   3. Configure your APIs in config/apis.json"
echo "   4. Set up webhooks for Slack/Discord notifications"
echo ""
echo "🔗 Useful links:"
echo "   • Railway Dashboard: https://railway.app/dashboard"
echo "   • Project Logs: railway logs"
echo "   • Environment Variables: railway variables"

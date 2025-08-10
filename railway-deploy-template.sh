#!/bin/bash

# Railway One-Click Deploy Template
# This script creates a Railway deployment template for the API Uptime Monitor

echo "🚀 Creating Railway Deployment Template"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    print_error "Railway CLI is not installed"
    echo "Please install it with: npm install -g @railway/cli"
    exit 1
fi

print_status "Railway CLI found"

# Check if user is logged in
if ! railway whoami &> /dev/null; then
    print_info "Logging into Railway..."
    railway login
    if ! railway whoami &> /dev/null; then
        print_error "Failed to login to Railway"
        exit 1
    fi
fi

print_status "Logged into Railway as $(railway whoami)"

# Create new Railway project
print_info "Creating new Railway project..."
if [ ! -f ".railway/project.json" ]; then
    railway init -n "API Uptime Monitor"
    print_status "Railway project created"
else
    print_warning "Railway project already exists"
fi

# Deploy MongoDB service first
print_info "Setting up MongoDB service..."
echo "Please follow these steps in your Railway dashboard:"
echo "1. Go to https://railway.app/dashboard"
echo "2. Open your 'API Uptime Monitor' project"
echo "3. Click 'New Service' → 'Database' → 'Add MongoDB'"
echo "4. Wait for MongoDB to deploy"
echo ""
read -p "Press Enter when MongoDB service is ready..."

# Get MongoDB connection details
print_info "Setting up environment variables..."

# Create environment variables
railway variables set PORT=8080
railway variables set GIN_MODE=release
railway variables set CHECK_INTERVAL="*/5 * * * *"
railway variables set TIMEOUT_SECONDS=30
railway variables set MAX_RETRIES=3
railway variables set DOWNTIME_THRESHOLD=3
railway variables set ENABLE_SLACK=false
railway variables set ENABLE_DISCORD=false
railway variables set DATABASE_NAME=uptime_monitor

print_status "Basic environment variables set"

# Deploy the application
print_info "Deploying the application..."
railway up --detach

print_status "Application deployed!"

# Get the deployment URL
RAILWAY_URL=$(railway domain 2>/dev/null || echo "Not available yet")

echo ""
echo "🎉 Deployment Complete!"
echo "======================="
echo ""
print_info "Your uptime monitor is being deployed to Railway"

if [ "$RAILWAY_URL" != "Not available yet" ]; then
    print_status "Access your dashboard at: https://$RAILWAY_URL"
else
    print_info "Dashboard URL will be available once deployment completes"
fi

echo ""
echo "📝 Next Steps:"
echo "=============="
echo "1. Wait for deployment to complete (check Railway dashboard)"
echo "2. Set up MongoDB connection string:"
echo "   - Go to Railway dashboard → MongoDB service → Connect tab"
echo "   - Copy the connection string"
echo "   - Set MONGODB_URI variable: railway variables set MONGODB_URI=\"<connection-string>\""
echo ""
echo "3. Configure your APIs to monitor:"
echo "   - Edit config/apis.json in your repository"
echo "   - Commit and push changes to trigger redeploy"
echo ""
echo "4. Set up notifications (optional):"
echo "   - Get Slack webhook URL and set: railway variables set SLACK_WEBHOOK_URL=\"<url>\""
echo "   - Enable Slack: railway variables set ENABLE_SLACK=true"
echo "   - Get Discord webhook URL and set: railway variables set DISCORD_WEBHOOK_URL=\"<url>\""
echo "   - Enable Discord: railway variables set ENABLE_DISCORD=true"
echo ""
echo "🔗 Useful Railway Commands:"
echo "=========================="
echo "• View logs: railway logs"
echo "• Check variables: railway variables"
echo "• Open dashboard: railway open"
echo "• Redeploy: railway up"
echo ""
print_status "Template deployment script completed!"

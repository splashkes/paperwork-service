#!/bin/bash

# Art Battle Paperwork Service - Deployment Script

set -e

echo "🚀 Deploying Art Battle Paperwork Service to DigitalOcean Apps..."

# Check if doctl is installed
if ! command -v doctl &> /dev/null; then
    echo "❌ doctl CLI is not installed. Please install it first."
    exit 1
fi

# Check if user is authenticated with doctl
if ! doctl auth whoami &> /dev/null; then
    echo "❌ Please authenticate with doctl first: doctl auth init"
    exit 1
fi

# Build the application locally to check for errors
echo "🔨 Building application locally..."
make build

# Test the build
if [ ! -f "bin/paperwork-service" ]; then
    echo "❌ Build failed - binary not found"
    exit 1
fi

echo "✅ Local build successful"

# Check if this is an initial deployment or update
APP_ID=""
if [ -f ".do/app_id" ]; then
    APP_ID=$(cat .do/app_id)
    echo "📱 Found existing app ID: $APP_ID"
else
    echo "🆕 This appears to be a new deployment"
fi

# Deploy to DigitalOcean
if [ -z "$APP_ID" ]; then
    echo "🚀 Creating new app on DigitalOcean..."

    # Create the app
    RESULT=$(doctl apps create .do/app.yaml --format ID --no-header)
    APP_ID=$RESULT

    # Save app ID for future deployments
    mkdir -p .do
    echo "$APP_ID" > .do/app_id

    echo "✅ App created with ID: $APP_ID"
else
    echo "🔄 Updating existing app..."

    # Update the app
    doctl apps update "$APP_ID" --spec .do/app.yaml

    echo "✅ App updated"
fi

echo "⏳ Triggering deployment..."
doctl apps create-deployment "$APP_ID" --force-rebuild

echo "📊 Getting app info..."
doctl apps get "$APP_ID"

echo ""
echo "🎉 Deployment initiated!"
echo "📱 App ID: $APP_ID"
echo "🌐 You can monitor the deployment at: https://cloud.digitalocean.com/apps/$APP_ID"
echo ""
echo "🔗 Once deployed, your service will be available at:"
echo "   https://your-app-url.ondigitalocean.app/api/v1/health"
echo ""
echo "📝 Don't forget to:"
echo "   1. Set environment variables in the DigitalOcean dashboard"
echo "   2. Configure custom domain if needed"
echo "   3. Test the paperwork generation endpoint"
echo ""
echo "🧪 Test command (replace with actual URL):"
echo "   curl 'https://your-app-url.ondigitalocean.app/api/v1/event-pdf/AB2995' -o test.pdf"
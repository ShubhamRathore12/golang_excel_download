#!/bin/bash

echo "🚀 Deploying Export API to Render..."

# Check if render CLI is installed
if ! command -v render &> /dev/null; then
    echo "❌ Render CLI is not installed. Please install it first:"
    echo "   npm install -g @render/cli"
    echo "   or visit: https://render.com/docs/cli"
    exit 1
fi

# Check if we're logged in to Render
if ! render whoami &> /dev/null; then
    echo "❌ Not logged in to Render. Please run: render login"
    exit 1
fi

echo "✅ Render CLI is ready"

# Create or update the service
echo "📦 Creating/updating Render service..."
render blueprint apply

if [ $? -eq 0 ]; then
    echo "✅ Deployment successful!"
    echo ""
    echo "🌐 Your API should be available at:"
    echo "   https://export-api.onrender.com"
    echo ""
    echo "📋 Available endpoints:"
    echo "   GET /health     - Health check"
    echo "   GET /tables     - List available tables"
    echo "   GET /status     - Service status"
    echo "   GET /export     - Export data to Excel"
    echo ""
    echo "🔧 To check deployment status:"
    echo "   render ps"
else
    echo "❌ Deployment failed. Check the error messages above."
    exit 1
fi

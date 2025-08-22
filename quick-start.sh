#!/bin/bash

echo "🚀 Quick Start for Export API"
echo "=============================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is installed, but version $REQUIRED_VERSION or higher is required."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Set environment variables
echo "🔧 Setting environment variables..."
export DB_HOST=myshaa.com
export DB_PORT=3306
export DB_USER=myshaa_kabu
export DB_PASSWORD=T-Cyj;f5g1y6
export DB_NAME=myshaa_kabu
export PORT=8080

echo "📦 Installing dependencies..."
go mod download

if [ $? -ne 0 ]; then
    echo "❌ Failed to download dependencies"
    exit 1
fi

echo "🔨 Building application..."
go build -o export-api main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful!"

echo ""
echo "🌐 Starting Export API server..."
echo "   Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo "   Port: $PORT"
echo "   User: $DB_USER"
echo ""
echo "📋 Available endpoints:"
echo "   GET /health     - Health check"
echo "   GET /tables     - List available tables"
echo "   GET /status     - Service status"
echo "   GET /export     - Export data to Excel"
echo ""
echo "🔍 Testing the API..."
echo "   Health check: curl http://localhost:$PORT/health"
echo "   Tables: curl http://localhost:$PORT/tables"
echo "   Status: curl http://localhost:$PORT/status"
echo ""
echo "📤 Example export:"
echo "   curl \"http://localhost:$PORT/export?table=GTPL_108_gT_40E_P_S7_200_Germany&all=true\" -o export.xlsx"
echo ""
echo "⏹️  Press Ctrl+C to stop the server"
echo ""

# Start the server
./export-api

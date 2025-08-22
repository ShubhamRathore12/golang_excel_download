#!/bin/bash

echo "Setting up Export API Go project..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "Error: Go version $GO_VERSION is installed, but version $REQUIRED_VERSION or higher is required."
    exit 1
fi

echo "Go version $GO_VERSION detected ✓"

# Initialize Go module
echo "Initializing Go module..."
go mod tidy

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Verify dependencies
echo "Verifying dependencies..."
go mod verify

# Build the application
echo "Building application..."
go build -o bin/export-api main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful! Binary created at bin/export-api"
    echo ""
    echo "Next steps:"
    echo "1. Set your database environment variables:"
    echo "   export DB_HOST=myshaa.com"
    echo "   export DB_PORT=3306"
    echo "   export DB_USER=myshaa_kabu"
    echo "   export DB_PASSWORD=T-Cyj;f5g1y6"
    echo "   export DB_NAME=myshaa_kabu"
    echo ""
    echo "2. Run the application:"
    echo "   ./bin/export-api"
    echo ""
    echo "3. Or run directly with:"
    echo "   go run main.go"
    echo ""
    echo "4. Test the API:"
    echo "   curl http://localhost:8080/health"
    echo "   curl http://localhost:8080/tables"
    echo "   curl http://localhost:8080/status"
else
    echo "✗ Build failed. Please check the error messages above."
    exit 1
fi

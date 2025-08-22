.PHONY: help build run test clean docker-build docker-run deploy

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  deploy       - Deploy to Render (requires render CLI)"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/export-api main.go
	@echo "Build complete: bin/export-api"

# Run the application locally
run:
	@echo "Running application..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf dist/
	go clean
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t export-api .
	@echo "Docker image built: export-api"

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 \
		-e DB_HOST=localhost \
		-e DB_PORT=3306 \
		-e DB_USER=root \
		-e DB_PASSWORD=password \
		-e DB_NAME=test \
		export-api

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	gosec ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Performance benchmark
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Coverage report
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Pre-commit checks
pre-commit: fmt lint test
	@echo "Pre-commit checks complete"

# Development setup
dev-setup: deps fmt lint test
	@echo "Development environment setup complete"

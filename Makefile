# Makefile for Khalif Backend
# Auto-start PostgreSQL & Redis before running the app

.PHONY: dev run start-services stop-services build test clean

# Default target
dev: start-services run

# Start PostgreSQL and Redis services
start-services:
	@echo "🔄 Starting PostgreSQL..."
	@sudo service postgresql start 2>/dev/null || echo "PostgreSQL already running or failed to start"
	@echo "🔄 Starting Redis..."
	@sudo service redis-server start 2>/dev/null || echo "Redis already running or failed to start"
	@echo "✅ Services started!"
	@sleep 1

# Stop services
stop-services:
	@echo "🛑 Stopping PostgreSQL..."
	@sudo service postgresql stop 2>/dev/null || true
	@echo "🛑 Stopping Redis..."
	@sudo service redis-server stop 2>/dev/null || true
	@echo "✅ Services stopped!"

# Run the application
run:
	@echo "🚀 Starting Khalif Backend..."
	go run cmd/api/main.go

# Build the application
build:
	@echo "📦 Building application..."
	go build -o bin/khalif-backend cmd/api/main.go
	@echo "✅ Build complete: bin/khalif-backend"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	@echo "✅ Clean complete!"

# Check service status
status:
	@echo "📊 Service Status:"
	@echo "PostgreSQL:"
	@sudo service postgresql status 2>/dev/null || echo "  Not running"
	@echo "Redis:"
	@sudo service redis-server status 2>/dev/null || echo "  Not running"

# Install dependencies
deps:
	@echo "📥 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

# Help
help:
	@echo "Available commands:"
	@echo "  make dev      - Start services and run the app (default)"
	@echo "  make run      - Run the app without starting services"
	@echo "  make start-services - Start PostgreSQL & Redis"
	@echo "  make stop-services  - Stop PostgreSQL & Redis"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make status   - Check service status"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make help     - Show this help"

# Makefile for Khalif Backend
# Auto-start PostgreSQL & Redis before running the app

.PHONY: dev run start-services stop-services build test clean migrate-up migrate-down migrate-status migrate-create

# Database configuration (loaded from .env or defaults)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= khalif_db
DB_SSLMODE ?= disable
DATABASE_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Default target
dev: start-services run

# Start PostgreSQL, Redis, and Meilisearch services
start-services:
	@echo "🔄 Starting PostgreSQL..."
	@sudo service postgresql start 2>/dev/null || echo "PostgreSQL already running or failed to start"
	@echo "🔄 Starting Redis..."
	@sudo service redis-server start 2>/dev/null || echo "Redis already running or failed to start"
	@echo "🔄 Starting Meilisearch..."
	@if pgrep -x "meilisearch" > /dev/null; then \
		echo "Meilisearch already running"; \
	else \
		./meilisearch --master-key="khalif_search_key" > meilisearch.log 2>&1 & \
		echo "Meilisearch started (logs: meilisearch.log)"; \
	fi
	@echo "✅ Services started!"
	@sleep 1

# Stop services
stop-services:
	@echo "🛑 Stopping PostgreSQL..."
	@sudo service postgresql stop 2>/dev/null || true
	@echo "🛑 Stopping Redis..."
	@sudo service redis-server stop 2>/dev/null || true
	@echo "🛑 Stopping Meilisearch..."
	@pkill -x meilisearch 2>/dev/null || true
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
	@echo "Meilisearch:"
	@pgrep -x meilisearch > /dev/null && echo "  Running" || echo "  Not running"

# Install dependencies
deps:
	@echo "📥 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

# ============================================
# Migration Commands (golang-migrate)
# ============================================

# Run all pending migrations
migrate-up:
	@echo "⬆️  Running migrations..."
	@migrate -path migrations/sql -database "$(DATABASE_URL)" up
	@echo "✅ Migrations applied!"

# Rollback last migration
migrate-down:
	@echo "⬇️  Rolling back last migration..."
	@migrate -path migrations/sql -database "$(DATABASE_URL)" down 1
	@echo "✅ Rollback complete!"

# Rollback all migrations
migrate-reset:
	@echo "🔄 Resetting all migrations..."
	@migrate -path migrations/sql -database "$(DATABASE_URL)" down -all
	@echo "✅ All migrations rolled back!"

# Show current migration version
migrate-version:
	@echo "📊 Current migration version:"
	@migrate -path migrations/sql -database "$(DATABASE_URL)" version

# Force set migration version (use with caution!)
migrate-force:
	@echo "⚠️  Force setting migration version to $(VERSION)..."
	@migrate -path migrations/sql -database "$(DATABASE_URL)" force $(VERSION)
	@echo "✅ Version forced!"

# Create new migration file
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "📝 Creating migration: $(NAME)..."
	@migrate create -ext sql -dir migrations/sql -seq $(NAME)
	@echo "✅ Migration files created!"

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "  📦 Application:"
	@echo "  make dev      - Start services and run the app (default)"
	@echo "  make run      - Run the app without starting services"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make deps     - Install dependencies"
	@echo "  make clean    - Clean build artifacts"
	@echo ""
	@echo "  🔌 Services:"
	@echo "  make start-services - Start PostgreSQL, Redis & Meilisearch"
	@echo "  make stop-services  - Stop PostgreSQL, Redis & Meilisearch"
	@echo "  make status   - Check service status"
	@echo ""
	@echo "  🗄️  Migrations:"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-reset   - Rollback all migrations"
	@echo "  make migrate-version - Show current migration version"
	@echo "  make migrate-create NAME=xxx - Create new migration"
	@echo ""
	@echo "  make help     - Show this help"

# Database Migration Makefile
# This Makefile provides convenient commands for managing database migrations

.PHONY: help migrate-up migrate-down migrate-create install-migrate-cli

# Default target
help: ## Show this help message
	@echo "Database Migration Commands:"
	@echo ""
	@echo "  make migrate-up              Apply all pending database migrations"
	@echo "  make migrate-down            Rollback the most recent migration"
	@echo "  make migrate-create name=X   Create new migration files with name X"
	@echo "  make install-migrate-cli     Install golang-migrate CLI tool"
	@echo "  make help                    Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-create name=add_users_table"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo ""
	@echo "Prerequisites:"
	@echo "  - golang-migrate CLI tool (run 'make install-migrate-cli' to install)"
	@echo "  - Database configuration in resources/config/config.yaml"
	@echo "  - Migration files in resources/migrations directory"

# Install golang-migrate CLI tool
install-migrate-cli: ## Install golang-migrate CLI tool
	@echo "Installing golang-migrate CLI..."
	@if command -v migrate >/dev/null 2>&1; then \
		echo "golang-migrate CLI is already installed"; \
		migrate -version; \
	else \
		echo "Installing golang-migrate CLI..."; \
		if command -v brew >/dev/null 2>&1; then \
			brew install golang-migrate; \
		elif command -v go >/dev/null 2>&1; then \
			go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		else \
			echo "Error: Neither brew nor go found. Please install golang-migrate manually."; \
			echo "Visit: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
			exit 1; \
		fi; \
	fi

# Check if migrate CLI is available
check-migrate-cli:
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "Error: golang-migrate CLI not found."; \
		echo "Please run 'make install-migrate-cli' first."; \
		exit 1; \
	fi

# Ensure migrations directory exists
ensure-migrations-dir:
	@mkdir -p resources/migrations

# Database URL generation helper
get-db-url:
	@go run scripts/get-db-url.go

# Apply all pending database migrations
migrate-up: check-migrate-cli ensure-migrations-dir ## Apply all pending database migrations
	@echo "Applying pending database migrations..."
	@DB_URL=$$(go run scripts/get-db-url.go 2>/dev/null); \
	if [ $$? -ne 0 ]; then \
		echo "Error: Failed to get database configuration"; \
		echo "Please ensure resources/config/config.yaml exists and is properly configured"; \
		exit 1; \
	fi; \
	if [ -z "$$DB_URL" ]; then \
		echo "Error: Database URL is empty"; \
		echo "Please check your database configuration in resources/config/config.yaml"; \
		exit 1; \
	fi; \
	echo "Connecting to database..."; \
	migrate -path resources/migrations -database "$$DB_URL" up; \
	MIGRATE_EXIT_CODE=$$?; \
	if [ $$MIGRATE_EXIT_CODE -eq 0 ]; then \
		echo "✅ Database migrations applied successfully"; \
	else \
		echo "❌ Migration failed with exit code $$MIGRATE_EXIT_CODE"; \
		echo "Please check the error messages above and verify:"; \
		echo "  - Database is running and accessible"; \
		echo "  - Database credentials are correct"; \
		echo "  - Migration files are valid SQL"; \
		exit $$MIGRATE_EXIT_CODE; \
	fi

# Rollback the most recent database migration
migrate-down: check-migrate-cli ensure-migrations-dir ## Rollback the most recent database migration
	@echo "Rolling back the most recent database migration..."
	@DB_URL=$$(go run scripts/get-db-url.go 2>/dev/null); \
	if [ $$? -ne 0 ]; then \
		echo "Error: Failed to get database configuration"; \
		echo "Please ensure resources/config/config.yaml exists and is properly configured"; \
		exit 1; \
	fi; \
	if [ -z "$$DB_URL" ]; then \
		echo "Error: Database URL is empty"; \
		echo "Please check your database configuration in resources/config/config.yaml"; \
		exit 1; \
	fi; \
	echo "Connecting to database..."; \
	migrate -path resources/migrations -database "$$DB_URL" down 1; \
	MIGRATE_EXIT_CODE=$$?; \
	if [ $$MIGRATE_EXIT_CODE -eq 0 ]; then \
		echo "✅ Database migration rolled back successfully"; \
	else \
		echo "❌ Migration rollback failed with exit code $$MIGRATE_EXIT_CODE"; \
		echo "Please check the error messages above and verify:"; \
		echo "  - Database is running and accessible"; \
		echo "  - Database credentials are correct"; \
		echo "  - There are migrations to rollback"; \
		echo "  - Down migration files are valid SQL"; \
		exit $$MIGRATE_EXIT_CODE; \
	fi

# Create new migration files with timestamped names
migrate-create: check-migrate-cli ensure-migrations-dir ## Create new migration files (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Migration name is required"; \
		echo ""; \
		echo "Usage:"; \
		echo "  make migrate-create name=<migration_name>"; \
		echo ""; \
		echo "Examples:"; \
		echo "  make migrate-create name=add_users_table"; \
		echo "  make migrate-create name=create_orders_index"; \
		echo "  make migrate-create name=update_user_schema"; \
		echo ""; \
		echo "Note: Migration name should be descriptive and use underscores for spaces"; \
		exit 1; \
	fi; \
	echo "Creating new migration files for: $(name)"; \
	migrate create -ext sql -dir resources/migrations -seq $(name); \
	CREATE_EXIT_CODE=$$?; \
	if [ $$CREATE_EXIT_CODE -eq 0 ]; then \
		UP_FILE=$$(ls resources/migrations/*$(name).up.sql 2>/dev/null | head -1); \
		DOWN_FILE=$$(ls resources/migrations/*$(name).down.sql 2>/dev/null | head -1); \
		if [ -n "$$UP_FILE" ] && [ -n "$$DOWN_FILE" ]; then \
			echo "-- Migration: $(name)" > "$$UP_FILE"; \
			echo "-- Created: $$(date)" >> "$$UP_FILE"; \
			echo "" >> "$$UP_FILE"; \
			echo "-- Add your SQL statements here" >> "$$UP_FILE"; \
			echo "-- Example:" >> "$$UP_FILE"; \
			echo "-- CREATE TABLE example_table (" >> "$$UP_FILE"; \
			echo "--     id SERIAL PRIMARY KEY," >> "$$UP_FILE"; \
			echo "--     name VARCHAR(255) NOT NULL," >> "$$UP_FILE"; \
			echo "--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP" >> "$$UP_FILE"; \
			echo "-- );" >> "$$UP_FILE"; \
			echo "" >> "$$UP_FILE"; \
			echo "-- Migration: $(name) (rollback)" > "$$DOWN_FILE"; \
			echo "-- Created: $$(date)" >> "$$DOWN_FILE"; \
			echo "" >> "$$DOWN_FILE"; \
			echo "-- Add your rollback SQL statements here" >> "$$DOWN_FILE"; \
			echo "-- Example:" >> "$$DOWN_FILE"; \
			echo "-- DROP TABLE IF EXISTS example_table;" >> "$$DOWN_FILE"; \
			echo "" >> "$$DOWN_FILE"; \
		fi; \
		echo "✅ Migration files created successfully in resources/migrations/"; \
		echo ""; \
		echo "Next steps:"; \
		echo "1. Edit the .up.sql file to add your schema changes"; \
		echo "2. Edit the .down.sql file to add the rollback logic"; \
		echo "3. Run 'make migrate-up' to apply the migration"; \
		echo ""; \
		echo "Created files:"; \
		ls -la resources/migrations/*$(name)* 2>/dev/null || echo "Files created with timestamp prefix"; \
	else \
		echo "❌ Failed to create migration files with exit code $$CREATE_EXIT_CODE"; \
		echo "Please check:"; \
		echo "  - golang-migrate CLI is properly installed"; \
		echo "  - resources/migrations directory is writable"; \
		echo "  - Migration name contains only valid characters"; \
		exit $$CREATE_EXIT_CODE; \
	fi
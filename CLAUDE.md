# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Running the Application
```bash
go run main.go
```

### Database Migrations
```bash
# Install migration tool (first time only)
make install-migrate-cli

# Apply all pending migrations
make migrate-up

# Rollback the most recent migration
make migrate-down

# Create new migration files
make migrate-create name=add_subscriptions_table
```

Note: The Makefile references `scripts/get-db-url.go` which needs to be created. This script should read the database configuration from `resources/config/config.yaml` and output a PostgreSQL connection string.

### Testing
```bash
# Run all tests (when implemented)
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./controllers/...
```

## Architecture Overview

This is the **Billing Service (jbilling)** for the Froconnect SaaS platform - an internal microservice responsible for all payment processing, subscription management, and billing operations.

### Service Position in Architecture
```
Contact Form Service → Billing Service → Stripe
                         ↓
                    PostgreSQL
```

The Billing Service is **never accessed directly by the frontend**. All requests come through the Contact Form Service.

### Key Architectural Decisions

1. **Configuration Management**: Uses Viper with hot-reload capability. Configuration files are loaded from multiple paths with 5-minute caching for performance. The `ConfigManager` is a thread-safe singleton.

2. **Database Pattern**: 
   - Connection pooling with pgx/v5
   - SQL migrations using golang-migrate
   - Migrations run automatically on startup
   - Each service has its own database

3. **Web Framework**: Chi v5 router with structured middleware chain (Logger → Recoverer → CORS)

4. **Logging**: Dual output to stdout and `/var/log/APP/jbilling/jbilling.log` using slog

### Code Organization

- `conns/` - External connections and infrastructure
  - `configs/` - Configuration management with hot-reload
  - `databases/` - Database connection and migration logic
- `controllers/` - HTTP handlers (to be implemented)
  - `dto/` - Standardized error responses
- `webapp/` - Server lifecycle management with graceful shutdown

### Database Schema

**Customers Table**: Links Stripe customers to organizations
- `organization_id` is the foreign key to the Contact Form Service
- `payment_method_id` references Stripe payment methods
- `default_payment_method` stores JSON data from Stripe

**Plans Table**: Defines subscription tiers
- `features` and `limits` are JSONB for flexibility
- `slug` values should match: `free`, `basic`, `premium`, `elite`, `enterprise`

### Billing Logic Implementation Notes

1. **Upgrades**: Immediate with proration - charge the difference right away
2. **Downgrades**: Delayed until period end - prevent abuse
3. **Overages**: Track usage, bill at period end at $0.01/email
4. **Grace Period**: 7 days for failed payments before suspension

### Missing Components to Implement

1. **API Controllers**: No endpoints implemented yet
2. **Stripe Integration**: Stripe SDK not added to go.mod
3. **Scripts**: `scripts/get-db-url.go` needs to be created
4. **Additional Tables**: Need subscriptions, usage_tracking, invoices tables
5. **Authentication**: JWT validation middleware
6. **Tests**: No test files exist yet

### Configuration

The service expects configuration at `resources/config/config.yaml`. In production, it also checks `/etc/APP/revonoir/jbilling/`.

Current configuration structure:
```yaml
database:
  host: localhost
  port: 5432
  user: jbilling
  password: <password>
  dbname: jbilling_db
  max_conns: 10
  sslmode: disable
```

### Development Workflow

1. Make database schema changes via migrations
2. Run `make migrate-up` to apply changes
3. Implement corresponding Go structs and logic
4. The server auto-reloads configuration changes without restart
5. Logs are structured JSON for easy parsing
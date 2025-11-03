# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `jrender`, a Go microservice for the Froconnect SaaS platform that serves embeddable contact forms. It operates as a public endpoint without API key authentication, relying on domain verification, rate limiting, and other security measures.

## Architecture

### Core Components

- **Configuration Management**: Uses Viper with YAML configs, supports live reloading with 5-minute cache TTL
- **Database Layer**: PostgreSQL with pgx/v5 driver and connection pooling
- **Web Framework**: Chi v5 router with CORS middleware
- **Migrations**: golang-migrate with sequential SQL files
- **Error Handling**: Custom `jerrors` package with structured JSON error responses

### Key Packages

- `conns/configs/`: Configuration management with caching and hot-reload
- `conns/databases/`: Database connection, pooling, and migration utilities
- `controllers/dto/jerrors/`: Standardized error response structures
- `webapp/`: HTTP server with graceful shutdown handling
- `scripts/`: Database URL generation for migrations

### Configuration

Config files located in:
- `/etc/APP/revonoir/jrender/config.yaml` (production)
- `./resources/config/config.yaml` (development)

Structure includes database settings and remote service URLs (jform service).

## Development Commands

### Database Migrations

All migration commands are handled through the Makefile:

```bash
# Apply all pending migrations
make migrate-up

# Rollback the most recent migration
make migrate-down

# Create new migration files
make migrate-create name=migration_name

# Install golang-migrate CLI
make install-migrate-cli
```

### Running the Application

```bash
# Start the server (listens on :9200)
go run main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Database Setup

```bash
# Create database
createdb jrender_db

# Copy config template
cp resources/config/config.yaml.example resources/config/config.yaml
# Edit with your database credentials

# Run migrations
make migrate-up
```

## Important Implementation Details

### Configuration Pattern
- `ConfigManager` implements thread-safe config loading with TTL caching
- Supports live config reloading via fsnotify
- Always use `configManager.GetConfig()` rather than direct Viper calls

### Database Pattern
- Use `NewDatabase(ctx, config)` for connection initialization
- Database URL generation handled by `scripts/get-db-url.go`
- Always use context-aware database operations

### Error Handling Pattern
```go
// Use jerrors for consistent API responses
return jerrors.BadRequest("invalid input")
jerrors.WriteErrorResponse(w, err)
```

### Server Architecture
- Server runs on port 9200 with graceful shutdown
- CORS configured for cross-origin embedding
- Middleware: Logger, Recoverer, CORS

### Security Architecture
The service implements public form embedding security through:
- Domain verification via `embed_registrations` table
- Rate limiting at HAProxy level
- Captcha integration (configurable)
- Secure cookie attributes
- HTTPS-only policy enforcement
- Frame protection headers (X-Frame-Options, CSP)

## Database Schema

Key table: `embed_registrations`
- `form_id`: References forms from main jform service
- `allowed_domains`: Array of whitelisted domains for embedding
- GIN index on `allowed_domains` for efficient domain queries

## Testing

No test files exist yet. When adding tests:
- Follow Go testing conventions with `_test.go` files
- Use table-driven tests for multiple scenarios
- Test database operations with testcontainers or similar

## Migration Workflow

1. Create migration: `make migrate-create name=descriptive_name`
2. Edit `.up.sql` and `.down.sql` files in `resources/migrations/`
3. Apply: `make migrate-up`
4. Test rollback: `make migrate-down` (in development only)

Migrations are automatically templated with metadata and examples.
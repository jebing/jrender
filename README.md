# Froconnect Embed Service

## Overview

The Froconnect Embed Service is a component of the Froconnect SaaS platform, responsible for rendering the form to the html, css, and javascript, so they can be displayed in the direct link, html/javascript, or iFrame

## Architecture

The Froconnect Embed Service operates as an internal microservice within the Froconnect ecosystem:

```
┌─────────────────┐
│       User      │
└────────┬────────┘
         │ Embed API
         ▼
┌─────────────────┐
│  Embed Service  │
│  (jrender)  │
└────────┬────────┘
         │ Form API
         ▼
┌─────────────────┐
│   Form Service  │
│     (jform)     │
└────────┬────────┘
```

### Key Responsibilities

- **Generate Html, CSS, Javascript**: Generate the Html, CSS, and Javascript so the form can be viewed by users
- **Manage Domain Permission**: Manage which domains are allowed to display the form
- **Forward Form Submission**: Forward form submission
- **Security**: Evaluate captcha, enforce HTTPS-only policy, implement frame protection headers, and input validation

## Tech Stack

- **Language**: Go 1.24.1
- **Database**: PostgreSQL
- **Web Framework**: Chi v5
- **Database Driver**: pgx/v5
- **Migration Tool**: golang-migrate
- **Configuration**: Viper

## Project Structure

```
jrender/
├── conns/                  # Connection management
│   ├── configs/           # Configuration handling
│   └── databases/         # Database connections
├── controllers/           # HTTP controllers
│   └── dto/              # Data transfer objects
├── internal/              # Internal packages
├── pkg/                   # Public packages
│   └── utils/            # Utility functions
├── resources/             # Resources
│   ├── config/           # Configuration files
│   └── migrations/       # Database migrations
├── webapp/                # Web application setup
├── main.go               # Application entry point
├── go.mod                # Go modules
└── Makefile              # Build automation
```

## Configuration

Configuration is managed via `resources/config/config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: jrender
  password: <password>
  dbname: jrender_db
  max_conns: 10
  sslmode: disable
```

## Development Setup

### Prerequisites
- Go 1.24.1 or higher
- PostgreSQL 14+
- Make (for build automation)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd jrender
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
createdb jrender_db
```

4. Configure the application:
```bash
cp resources/config/config.yaml.example resources/config/config.yaml
# Edit config.yaml with your database credentials
```

5. Run migrations:
```bash
make migrate-up
```

6. Start the service:
```bash
go run main.go
```

## Database Migrations

Migrations are managed using golang-migrate:

```bash
# Create a new migration
make migrate-create name=<migration_name>

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
make migrate-status
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./controllers/...
```


## Monitoring & Logging

- Structured logging with slog
- Log files stored in `/var/log/APP/jrender/`
- Metrics exposed for Prometheus scraping
- Health check endpoint at `/health`


## Security Features

The service implements multiple layers of security for public form embedding:

### Core Security Measures
- **Domain Verification**: Whitelist-based domain validation for authorized embedding
- **Rate Limiting**: Configurable request throttling at HAProxy level
- **Captcha Integration**: User-configurable bot protection
- **Secure Cookies**: Cross-site and secure cookie attributes

### Enhanced Security Features
- **HTTPS-Only Policy**: Reject HTTP origins to ensure encrypted connections
- **Frame Protection**: X-Frame-Options and CSP headers to prevent unauthorized framing
- **Input Validation**: Comprehensive sanitization of form submission data
- **CORS Configuration**: Strict cross-origin resource sharing policies

## Support

For issues and questions:
- Create an issue in the project repository
- Contact the development team via Linear
- Check the Froconnect documentation

## License

This is a proprietary component of the Froconnect platform.
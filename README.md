# Froconnect Billing Service

## Overview

The Billing Service is a critical component of the Froconnect SaaS platform, responsible for managing all payment processing, subscription management, and billing operations. This service is built with Go and integrates with Stripe for payment processing.

## Architecture

The Billing Service operates as an internal microservice within the Froconnect ecosystem:

```
┌─────────────────┐
│  Contact Form   │
│     Service     │
└────────┬────────┘
         │ Internal API
         ▼
┌─────────────────┐
│ Billing Service │
│   (jbilling)    │
└────────┬────────┘
         │ Stripe API
         ▼
    ┌─────────┐
    │ Stripe  │
    └─────────┘
```

### Key Responsibilities

- **Customer Management**: Creating and managing Stripe customers
- **Subscription Handling**: Creating, updating, and canceling subscriptions
- **Payment Processing**: Managing payment methods and processing charges
- **Invoice Management**: Generating and tracking invoices
- **Usage Tracking**: Recording and billing for email usage
- **Billing Webhooks**: Processing Stripe events
- **Proration Logic**: Handling upgrade/downgrade calculations

## Tech Stack

- **Language**: Go 1.24.1
- **Database**: PostgreSQL
- **Web Framework**: Chi v5
- **Database Driver**: pgx/v5
- **Migration Tool**: golang-migrate
- **Configuration**: Viper
- **Payment Gateway**: Stripe (to be integrated)

## Project Structure

```
jbilling/
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

## Database Schema

### Customers Table
Stores Stripe customer information linked to organizations:

```sql
- id: UUID (Primary Key)
- organization_id: UUID (Unique, Not Null)
- email: VARCHAR(255)
- name: VARCHAR(255)
- billing_email: VARCHAR(255)
- payment_method_id: VARCHAR(255)
- default_payment_method: JSONB
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

### Plans Table
Defines available subscription plans:

```sql
- id: UUID (Primary Key)
- name: VARCHAR(100)
- slug: VARCHAR(100) (Unique)
- description: TEXT
- price_monthly: DECIMAL(10,2)
- price_yearly: DECIMAL(10,2)
- features: JSONB
- limits: JSONB
- is_active: BOOLEAN
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

## Pricing Structure

| Plan | Monthly | Yearly | Included Emails | Extra Email Price | Max Emails |
|------|---------|---------|-----------------|-------------------|------------|
| Free | $0 | $0 | 25 | N/A | 25 |
| Basic | $0.99 | $9.99 | 100 | $0.01 | 200 |
| Premium | $4.99 | $49.99 | 1,000 | $0.01 | 2,000 |
| Elite | $9.99 | $99.99 | 3,000 | $0.01 | 6,000 |
| Enterprise | Custom | Custom | Custom | Custom | Custom |

## API Endpoints

### Customer Management
- `POST /api/v1/customers` - Create a new customer
- `GET /api/v1/customers/:organization_id` - Get customer by organization
- `PUT /api/v1/customers/:organization_id` - Update customer
- `DELETE /api/v1/customers/:organization_id` - Delete customer

### Subscription Management
- `POST /api/v1/subscriptions` - Create subscription
- `GET /api/v1/subscriptions/:organization_id` - Get active subscription
- `PUT /api/v1/subscriptions/:organization_id` - Update subscription
- `DELETE /api/v1/subscriptions/:organization_id` - Cancel subscription

### Payment Methods
- `POST /api/v1/payment-methods` - Add payment method
- `GET /api/v1/payment-methods/:organization_id` - List payment methods
- `PUT /api/v1/payment-methods/:id/default` - Set default payment method
- `DELETE /api/v1/payment-methods/:id` - Remove payment method

### Usage & Billing
- `POST /api/v1/usage` - Record email usage
- `GET /api/v1/usage/:organization_id` - Get usage statistics
- `GET /api/v1/invoices/:organization_id` - List invoices
- `GET /api/v1/invoices/:id/download` - Download invoice

### Webhooks
- `POST /api/v1/webhooks/stripe` - Stripe webhook endpoint

## Billing Logic

### Subscription Upgrades
- Immediate effect with usage-based proration
- Customer charged for the difference immediately
- Unused portion of current plan credited

### Subscription Downgrades
- Take effect at the end of the current billing period
- Prevents abuse and revenue loss
- Customer continues with current plan features until period end

### Overage Handling
- Automatically billed at $0.01 per email over limit
- Charged at the end of billing period
- No service interruption for overages

### Payment Failures
- 7-day grace period for failed payments
- Automatic retry schedule: Day 1, 3, 5, 7
- Service suspension after grace period
- Reactivation upon successful payment

## Configuration

Configuration is managed via `resources/config/config.yaml`:

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

## Development Setup

### Prerequisites
- Go 1.24.1 or higher
- PostgreSQL 14+
- Make (for build automation)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd jbilling
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
createdb jbilling_db
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

## Deployment

The service is designed to run as a containerized application:

```bash
# Build Docker image
docker build -t froconnect-billing .

# Run container
docker run -p 8080:8080 --env-file .env froconnect-billing
```

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string
- `STRIPE_SECRET_KEY` - Stripe API secret key
- `STRIPE_WEBHOOK_SECRET` - Stripe webhook signing secret
- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Logging level (info, debug, error)

## Security Considerations

- All API endpoints require authentication via JWT tokens
- Stripe webhook endpoints verify signatures
- Sensitive data (payment methods) stored only in Stripe
- Database credentials encrypted at rest
- TLS required for all external communications

## Monitoring & Logging

- Structured logging with slog
- Log files stored in `/var/log/APP/jbilling/`
- Metrics exposed for Prometheus scraping
- Health check endpoint at `/health`

## Future Enhancements

- [ ] Support for multiple payment providers
- [ ] Advanced analytics and reporting
- [ ] Automated dunning management
- [ ] Tax calculation integration
- [ ] Multi-currency support
- [ ] Revenue recognition features

## Support

For issues and questions:
- Create an issue in the project repository
- Contact the development team via Linear
- Check the Froconnect documentation

## License

This is a proprietary component of the Froconnect platform.
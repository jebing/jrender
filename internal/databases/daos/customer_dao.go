package daos

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"revonoir.com/jbilling/controllers/dto/jerrors"
	"revonoir.com/jbilling/internal/databases/models"
)

// Table and column constants
const (
	CUSTOMERS_TABLE                   = "customers"
	CUSTOMER_ID                       = "id"
	CUSTOMER_ORGANIZATION_ID          = "organization_id"
	CUSTOMER_EMAIL                    = "email"
	CUSTOMER_NAME                     = "name"
	CUSTOMER_BILLING_EMAIL            = "billing_email"
	CUSTOMER_DEFAULT_PAYMENT_PROVIDER = "default_payment_provider"
	CUSTOMER_CREATED_AT               = "created_at"
	CUSTOMER_UPDATED_AT               = "updated_at"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type CustomerDAO interface {
	Create(ctx context.Context, req models.CreateCustomerRequest) (*models.Customer, error)
	GetByOrganizationID(ctx context.Context, organizationID uuid.UUID) (*models.Customer, error)
}

type customerDAO struct {
	db *pgxpool.Pool
}

func NewCustomerDAO(db *pgxpool.Pool) CustomerDAO {
	return &customerDAO{
		db: db,
	}
}

func (d *customerDAO) Create(ctx context.Context, req models.CreateCustomerRequest) (*models.Customer, error) {
	query := psql.Insert(CUSTOMERS_TABLE).
		Columns(
			CUSTOMER_ORGANIZATION_ID,
			CUSTOMER_EMAIL,
			CUSTOMER_NAME,
			CUSTOMER_BILLING_EMAIL,
			CUSTOMER_DEFAULT_PAYMENT_PROVIDER,
		).
		Values(
			req.OrganizationID,
			req.Email,
			req.Name,
			req.BillingEmail,
			req.DefaultPaymentProvider,
		).
		Suffix("RETURNING " + CUSTOMER_ID)

	sql, args, err := query.ToSql()
	if err != nil {
		slog.Error("Failed to build insert query", "error", err)
		return nil, jerrors.NewErrorResp(http.StatusInternalServerError, "failed to build insert query")
	}

	var customer models.Customer
	err = d.db.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
	)

	if err != nil {
		slog.Error("Failed to create customer", "error", err, "organization_id", req.OrganizationID)
		return nil, jerrors.NewErrorResp(http.StatusInternalServerError, "failed to create customer")
	}

	customer.OrganizationID = req.OrganizationID
	customer.Email = req.Email
	customer.Name = req.Name
	customer.BillingEmail = req.BillingEmail
	customer.DefaultPaymentProvider = req.DefaultPaymentProvider

	slog.Info("Customer created successfully", "customer_id", customer.ID, "organization_id", customer.OrganizationID)
	return &customer, nil
}

func (d *customerDAO) GetByOrganizationID(ctx context.Context, organizationID uuid.UUID) (*models.Customer, error) {
	query := psql.Select(
		CUSTOMER_ID,
		CUSTOMER_ORGANIZATION_ID,
		CUSTOMER_EMAIL,
		CUSTOMER_NAME,
		CUSTOMER_BILLING_EMAIL,
		CUSTOMER_DEFAULT_PAYMENT_PROVIDER,
	).
		From(CUSTOMERS_TABLE).
		Where(squirrel.Eq{CUSTOMER_ORGANIZATION_ID: organizationID})

	sql, args, err := query.ToSql()
	if err != nil {
		slog.Error("Failed to build select query", "error", err)
		return nil, err
	}

	var customer models.Customer
	err = d.db.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
		&customer.OrganizationID,
		&customer.Email,
		&customer.Name,
		&customer.BillingEmail,
		&customer.DefaultPaymentProvider,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info("Customer not found", "organization_id", organizationID)
			return nil, jerrors.NewErrorResp(http.StatusNotFound, "customer not found")
		}
		slog.Error("Failed to get customer", "error", err, "organization_id", organizationID)
		return nil, jerrors.NewErrorResp(http.StatusInternalServerError, "failed to get customer")
	}

	slog.Info("Customer retrieved successfully", "customer_id", customer.ID, "organization_id", organizationID)
	return &customer, nil
}

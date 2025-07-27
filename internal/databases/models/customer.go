package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	OrganizationID           uuid.UUID  `json:"organization_id" db:"organization_id"`
	Email                    string     `json:"email" db:"email"`
	Name                     *string    `json:"name,omitempty" db:"name"`
	BillingEmail             *string    `json:"billing_email,omitempty" db:"billing_email"`
	DefaultPaymentProvider   *string    `json:"default_payment_provider,omitempty" db:"default_payment_provider"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateCustomerRequest struct {
	OrganizationID         uuid.UUID `json:"organization_id" validate:"required"`
	Email                  string    `json:"email" validate:"required,email"`
	Name                   *string   `json:"name,omitempty"`
	BillingEmail           *string   `json:"billing_email,omitempty" validate:"omitempty,email"`
	DefaultPaymentProvider *string   `json:"default_payment_provider,omitempty"`
}
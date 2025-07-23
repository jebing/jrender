-- Migration: create_customer_table
-- Created: Wed Jul 23 12:50:12 +08 2025

CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID UNIQUE NOT NULL,
    email TEXT NOT NULL,
    name TEXT,
    billing_email TEXT,
    default_payment_provider TEXT, -- 'stripe', 'paypal', 'razorpay', etc.
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_customers_organization_id ON customers (organization_id);


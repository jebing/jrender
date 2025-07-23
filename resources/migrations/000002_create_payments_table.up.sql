-- Migration: create_payments_table
-- Created: Wed Jul 23 14:21:44 +08 2025

CREATE TABLE payment_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL, -- 'stripe', 'paypal', 'razorpay'
    display_name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    configuration JSONB, -- Provider-specific config (API endpoints, etc.)
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Seed the payment providers table with the default payment providers
INSERT INTO payment_providers (name, display_name, is_active, configuration) VALUES
  (
      'stripe',
      'Stripe',
      TRUE,
      '{
          "api_base_url": "https://api.stripe.com",
          "webhook_endpoint_secret": null,
          "supported_currencies": ["usd", "eur", "gbp"],
          "supported_payment_methods": ["card", "ach", "sepa_debit"],
          "features": {
              "subscriptions": true,
              "one_time_payments": true,
              "refunds": true,
              "webhooks": true,
              "proration": true
          }
      }'::jsonb
  );

CREATE TABLE customer_payment_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES payment_providers(id),
    provider_customer_id TEXT NOT NULL, -- stripe_customer_id, paypal_payer_id, etc.
    provider_customer_data JSONB, -- Provider-specific customer data
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(customer_id, provider_id),
    UNIQUE(provider_id, provider_customer_id)
);

CREATE INDEX idx_customer_payment_profiles_customer_id ON customer_payment_profiles (customer_id);
CREATE INDEX idx_customer_payment_profiles_provider_id ON customer_payment_profiles (provider_id);

CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    payment_profile_id UUID NOT NULL REFERENCES customer_payment_profiles(id) ON DELETE CASCADE,
    provider_payment_method_id VARCHAR(255) NOT NULL, -- provider's ID for the method
    type TEXT NOT NULL, -- 'card', 'bank_account', 'upi', 'wallet'
    display_name TEXT, -- '•••• 4242', 'john@upi', etc.
    metadata JSONB, -- last4, brand, exp_month, etc.
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
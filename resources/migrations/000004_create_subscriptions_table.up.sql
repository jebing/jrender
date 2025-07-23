-- Migration: create_subscriptions_table
-- Created: Wed Jul 23 14:26:15 +08 2025

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id),
    plan_id UUID NOT NULL REFERENCES plans(id),
    payment_profile_id UUID NOT NULL REFERENCES customer_payment_profiles(id),
    provider_subscription_id TEXT, -- Provider's subscription ID
    status TEXT NOT NULL, -- 'active', 'canceled', 'past_due', etc.
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    cancel_at_period_end BOOLEAN DEFAULT FALSE,
    metadata JSONB, -- Provider-specific subscription data
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_customer_id ON subscriptions (customer_id);
CREATE INDEX idx_subscriptions_status ON subscriptions (status);


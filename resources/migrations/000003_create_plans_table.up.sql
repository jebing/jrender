-- Migration: create_plans_table
-- Created: Wed Jul 23 14:25:32 +08 2025

CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    description TEXT,
    price_monthly DECIMAL(10,2),
    price_yearly DECIMAL(10,2),
    features JSONB NOT NULL DEFAULT '{}',
    limits JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_plans_slug ON plans (slug);

-- Seed the plans table with the default plans
INSERT INTO plans (name, slug, description, price_monthly, price_yearly, features, limits, is_active) VALUES
  (
      'Free',
      'free',
      'Perfect for trying out Froconnect with basic contact form functionality.',
      0.00,
      0.00,
      '{
          "forms": "unlimited",
          "templates": "basic",
          "custom_branding": false,
          "custom_domain": false,
          "analytics": "basic",
          "spam_protection": true,
          "email_notifications": true,
          "api_access": false,
          "support": "community"
      }'::jsonb,
      '{
          "emails_per_month": 25,
          "max_forms": 1,
          "max_fields_per_form": 10,
          "file_uploads": false,
          "storage_mb": 0,
          "team_members": 1
      }'::jsonb,
      TRUE
  ),
  (
      'Basic',
      'basic',
      'Great for small businesses and personal projects with moderate email volume.',
      0.99,
      9.99,
      '{
          "forms": "unlimited",
          "templates": "standard",
          "custom_branding": true,
          "custom_domain": false,
          "analytics": "standard",
          "spam_protection": true,
          "email_notifications": true,
          "api_access": true,
          "support": "email"
      }'::jsonb,
      '{
          "emails_per_month": 100,
          "max_emails_total": 200,
          "overage_price_per_email": 0.01,
          "max_forms": 5,
          "max_fields_per_form": 20,
          "file_uploads": true,
          "storage_mb": 100,
          "team_members": 2
      }'::jsonb,
      TRUE
  ),
  (
      'Premium',
      'premium',
      'Perfect for growing businesses with higher email volumes and advanced features.',
      4.99,
      49.99,
      '{
          "forms": "unlimited",
          "templates": "premium",
          "custom_branding": true,
          "custom_domain": true,
          "analytics": "advanced",
          "spam_protection": true,
          "email_notifications": true,
          "api_access": true,
          "support": "priority_email",
          "webhooks": true,
          "multi_language": true
      }'::jsonb,
      '{
          "emails_per_month": 1000,
          "max_emails_total": 2000,
          "overage_price_per_email": 0.01,
          "max_forms": 25,
          "max_fields_per_form": 50,
          "file_uploads": true,
          "storage_mb": 1000,
          "team_members": 5
      }'::jsonb,
      TRUE
  ),
  (
      'Elite',
      'elite',
      'For established businesses requiring high email volumes and premium support.',
      9.99,
      99.99,
      '{
          "forms": "unlimited",
          "templates": "elite",
          "custom_branding": true,
          "custom_domain": true,
          "analytics": "advanced",
          "spam_protection": true,
          "email_notifications": true,
          "api_access": true,
          "support": "phone_and_email",
          "webhooks": true,
          "multi_language": true,
          "white_label": true,
          "advanced_integrations": true
      }'::jsonb,
      '{
          "emails_per_month": 3000,
          "max_emails_total": 6000,
          "overage_price_per_email": 0.01,
          "max_forms": 100,
          "max_fields_per_form": 100,
          "file_uploads": true,
          "storage_mb": 5000,
          "team_members": 15
      }'::jsonb,
      TRUE
  ),
  (
      'Enterprise',
      'enterprise',
      'Custom solution for large organizations with specific requirements and unlimited usage.',
      NULL,
      NULL,
      '{
          "forms": "unlimited",
          "templates": "custom",
          "custom_branding": true,
          "custom_domain": true,
          "analytics": "enterprise",
          "spam_protection": true,
          "email_notifications": true,
          "api_access": true,
          "support": "dedicated_account_manager",
          "webhooks": true,
          "multi_language": true,
          "white_label": true,
          "advanced_integrations": true,
          "sso": true,
          "custom_development": true,
          "sla": true
      }'::jsonb,
      '{
          "emails_per_month": "unlimited",
          "max_emails_total": "unlimited",
          "overage_price_per_email": "custom",
          "max_forms": "unlimited",
          "max_fields_per_form": "unlimited",
          "file_uploads": true,
          "storage_mb": "unlimited",
          "team_members": "unlimited"
      }'::jsonb,
      TRUE
  );
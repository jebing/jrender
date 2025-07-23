-- Migration: create_payments_table (rollback)
-- Created: Wed Jul 23 14:21:44 +08 2025

DROP TABLE IF EXISTS customer_payment_profiles;
DROP TABLE IF EXISTS payment_providers;


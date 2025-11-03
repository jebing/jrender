-- Migration: create_embed_registrations_table (rollback)
-- Created: Fri Aug 15 13:12:16 +08 2025

-- Drop the embed_registrations table and all related indexes
DROP INDEX IF EXISTS idx_embed_registrations_allowed_domains;
DROP INDEX IF EXISTS idx_embed_registrations_is_active;
DROP INDEX IF EXISTS idx_embed_registrations_form_id;
DROP TABLE IF EXISTS embed_registrations;


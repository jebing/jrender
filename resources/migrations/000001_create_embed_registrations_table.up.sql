-- Migration: create_embed_registrations_table
-- Created: Fri Aug 15 13:12:16 +08 2025

-- Create embed_registrations table for form embedding with domain validation
CREATE TABLE embed_registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL,
    allowed_domains TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_embed_registrations_form_id ON embed_registrations(form_id);
CREATE INDEX idx_embed_registrations_is_active ON embed_registrations(is_active);

-- GIN index for efficient domain array queries
CREATE INDEX idx_embed_registrations_allowed_domains ON embed_registrations USING GIN(allowed_domains);


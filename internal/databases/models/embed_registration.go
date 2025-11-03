package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// EmbedRegistration represents a form embedding registration with domain validation
type EmbedRegistration struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	FormID         uuid.UUID      `json:"form_id" db:"form_id"`
	AllowedDomains pq.StringArray `json:"allowed_domains" db:"allowed_domains"`
	IsActive       bool           `json:"is_active" db:"is_active"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

// NewEmbedRegistration creates a new embed registration instance
func NewEmbedRegistration(formID uuid.UUID, allowedDomains []string) *EmbedRegistration {
	now := time.Now()
	return &EmbedRegistration{
		ID:             uuid.New(),
		FormID:         formID,
		AllowedDomains: pq.StringArray(allowedDomains),
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// ContainsDomain checks if the given domain is in the allowed domains list
func (er *EmbedRegistration) ContainsDomain(domain string) bool {
	for _, allowedDomain := range er.AllowedDomains {
		if allowedDomain == domain {
			return true
		}
	}
	return false
}

// AddDomain adds a domain to the allowed domains list if not already present
func (er *EmbedRegistration) AddDomain(domain string) {
	if !er.ContainsDomain(domain) {
		er.AllowedDomains = append(er.AllowedDomains, domain)
		er.UpdatedAt = time.Now()
	}
}

// RemoveDomain removes a domain from the allowed domains list
func (er *EmbedRegistration) RemoveDomain(domain string) {
	for i, allowedDomain := range er.AllowedDomains {
		if allowedDomain == domain {
			er.AllowedDomains = append(er.AllowedDomains[:i], er.AllowedDomains[i+1:]...)
			er.UpdatedAt = time.Now()
			break
		}
	}
}

// SetDomains replaces the entire allowed domains list
func (er *EmbedRegistration) SetDomains(domains []string) {
	er.AllowedDomains = pq.StringArray(domains)
	er.UpdatedAt = time.Now()
}

// Activate marks the embed registration as active
func (er *EmbedRegistration) Activate() {
	er.IsActive = true
	er.UpdatedAt = time.Now()
}

// Deactivate marks the embed registration as inactive
func (er *EmbedRegistration) Deactivate() {
	er.IsActive = false
	er.UpdatedAt = time.Now()
}

// GetDomainsSlice returns the allowed domains as a regular string slice
func (er *EmbedRegistration) GetDomainsSlice() []string {
	return []string(er.AllowedDomains)
}
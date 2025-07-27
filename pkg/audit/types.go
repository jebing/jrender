package audit

// Audit action constants for SMTP configuration operations
const (
	// Resource types
	ResourceUser   = "user"
	ResourceAPIKey = "api_key"

	// Event types
	EventTypeAudit = "audit"
)

// Context keys for storing audit information in Chi request context
type contextKey string

const (
	ClientIPKey  contextKey = "audit_client_ip"
	UserAgentKey contextKey = "audit_user_agent"
	UserIDKey    contextKey = "audit_user_id"
	RequestIDKey contextKey = "audit_request_id"
)

// AuditDetails represents additional context information for audit events
type AuditDetails map[string]interface{}

// NewAuditDetails creates a new AuditDetails map
func NewAuditDetails() AuditDetails {
	return make(map[string]interface{})
}

// Add adds a key-value pair to audit details
func (ad AuditDetails) Add(key string, value interface{}) AuditDetails {
	ad[key] = value
	return ad
}

// AddChange adds a field change to audit details (for update operations)
func (ad AuditDetails) AddChange(field, oldValue, newValue interface{}) AuditDetails {
	if ad["changes"] == nil {
		ad["changes"] = make(map[string]map[string]interface{})
	}
	changes := ad["changes"].(map[string]map[string]interface{})
	changes[field.(string)] = map[string]interface{}{
		"old": oldValue,
		"new": newValue,
	}
	return ad
}

// SanitizeValue removes sensitive information from values before logging
func SanitizeValue(key string, value interface{}) interface{} {
	sensitiveFields := map[string]bool{
		"password":      true,
		"smtp_password": true,
		"api_key":       true,
		"token":         true,
		"secret":        true,
		"credential":    true,
	}

	if sensitiveFields[key] {
		return "[REDACTED]"
	}

	// If value is a string and looks like a password (length check)
	if str, ok := value.(string); ok && len(str) > 8 &&
		(key == "smtp_password" || key == "password") {
		return "[REDACTED]"
	}

	return value
}

package audit

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// AuditMiddleware extracts request information and stores it in the Chi request context
// for use in audit logging throughout the request lifecycle
func AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client IP address
		clientIP := extractClientIP(r)

		// Extract User-Agent
		userAgent := r.Header.Get("User-Agent")

		// Get or generate request ID (Chi middleware.RequestID should be used before this)
		requestID := middleware.GetReqID(r.Context())
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Store audit information in Chi context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ClientIPKey, clientIP)
		ctx = context.WithValue(ctx, UserAgentKey, userAgent)
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateRequestID generates a simple request ID if Chi's RequestID middleware isn't used
func generateRequestID() string {
	// Simple fallback ID generation
	return "req_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:16]
}

// extractClientIP attempts to extract the real client IP address from various headers
func extractClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (most common proxy header)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header (used by some proxies)
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Check CF-Connecting-IP header (Cloudflare)
	cfIP := r.Header.Get("CF-Connecting-IP")
	if cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex]
	}

	return ip
}

// SetUserIDInContext sets the user ID in the Chi request context for audit logging
func SetUserIDInContext(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, userID)
	return r.WithContext(ctx)
}

// GetClientIP retrieves the client IP from the Chi request context
func GetClientIP(r *http.Request) string {
	if ip := r.Context().Value(ClientIPKey); ip != nil {
		if ipStr, ok := ip.(string); ok {
			return ipStr
		}
	}
	// Fallback to extracting directly from request
	return extractClientIP(r)
}

// GetUserAgent retrieves the user agent from the Chi request context
func GetUserAgent(r *http.Request) string {
	if ua := r.Context().Value(UserAgentKey); ua != nil {
		if uaStr, ok := ua.(string); ok {
			return uaStr
		}
	}
	// Fallback to extracting directly from request
	return r.Header.Get("User-Agent")
}

// GetUserID retrieves the user ID from the Chi request context
func GetUserID(r *http.Request) string {
	if userID := r.Context().Value(UserIDKey); userID != nil {
		if userIDStr, ok := userID.(string); ok {
			return userIDStr
		}
	}
	return ""
}

// GetRequestID retrieves the request ID from the Chi request context
func GetRequestID(r *http.Request) string {
	// First try to get from Chi's middleware
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		return reqID
	}

	// Fall back to our stored request ID
	if reqID := r.Context().Value(RequestIDKey); reqID != nil {
		if reqIDStr, ok := reqID.(string); ok {
			return reqIDStr
		}
	}
	return ""
}

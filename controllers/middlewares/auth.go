package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"revonoir.com/jbilling/controllers/dto"
	"revonoir.com/jbilling/controllers/dto/jerrors"
	"revonoir.com/jbilling/internal/databases/models"
	"revonoir.com/jbilling/pkg/audit"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
)

type ApiKeyAuthMiddlewareIf interface {
	AuthFunc() func(next http.Handler) http.Handler
}

type ApiKeyDaoIf interface {
	Get(ctx context.Context, apiKey string) (models.ApiKeys, error)
}

type ApiKeyAuthMiddleware struct {
	apiKeyDao ApiKeyDaoIf
}

func NewApiKeyAuthMiddleware(apiKeyDao ApiKeyDaoIf) *ApiKeyAuthMiddleware {
	return &ApiKeyAuthMiddleware{apiKeyDao: apiKeyDao}
}

func (m ApiKeyAuthMiddleware) AuthFunc() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := m.extractAPIKey(r)
			if err != nil {
				resp := dto.Response[any]{Error: &jerrors.ErrorResp{Message: err.Error()}}
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp)
				return
			}

			userId, err := m.validateAPIKey(r.Context(), tokenString)
			if err != nil {
				resp := dto.Response[any]{Error: &jerrors.ErrorResp{Message: err.Error()}}
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp)
				return
			}
			ctx := context.WithValue(r.Context(), audit.UserIDKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractAPIKey attempts to parse the API key from headers.
// Common patterns:
// - "Authorization: Bearer <token>"
// - "X-API-KEY: <token>"
func (m ApiKeyAuthMiddleware) extractAPIKey(r *http.Request) (string, error) {
	// Example: Checking for Bearer token in the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Maybe fallback to a custom header
		apiKeyHeader := r.Header.Get("X-API-KEY")
		if apiKeyHeader == "" {
			return "", jerrors.NewErrorResp(http.StatusUnauthorized, "no API key provided")
		}
		return apiKeyHeader, nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", jerrors.NewErrorResp(http.StatusUnauthorized, "invalid authorization header format")
	}

	return parts[1], nil
}

// validateAPIKey checks if the provided API key is valid.
// Typically this involves:
// 1. Hashing the incoming key.
// 2. Querying the database for a matching, active, non-expired key.
// 3. Returning associated user or service account ID if successful.
func (m ApiKeyAuthMiddleware) validateAPIKey(ctx context.Context, apiKey string) (string, error) {

	model, err := m.apiKeyDao.Get(ctx, apiKey)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", jerrors.NewErrorResp(http.StatusUnauthorized, "auth key is invalid")
		}
		slog.Error("Failed to get api key", "error", err)
		return "", jerrors.NewErrorResp(http.StatusInternalServerError, "Failed to get api key")
	}

	// check if the key is expired or revoked
	// Check expiration
	if model.ExpiresAt != nil && time.Now().After(*model.ExpiresAt) {
		slog.Error("Key expired", "apiKey", apiKey)
		return "", jerrors.NewErrorResp(http.StatusUnauthorized, "key expired")
	}
	// Check revocation
	if model.RevokedAt != nil {
		slog.Error("Key revoked", "apiKey", apiKey)
		return "", jerrors.NewErrorResp(http.StatusUnauthorized, "key revoked")
	}
	// Check status if necessary
	if model.Status != int8(models.ApiKeysStatusActive) {
		slog.Error("Key not active", "apiKey", apiKey)
		return "", jerrors.NewErrorResp(http.StatusUnauthorized, "key not active")
	}

	return model.UserID, nil
}

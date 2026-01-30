package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dusanbre/otg-sports-api/internal/api/auth"
	"github.com/dusanbre/otg-sports-api/internal/database"
)

// ContextKey type for context keys
type ContextKey string

const (
	// APIKeyContextKey is the context key for the API key
	APIKeyContextKey ContextKey = "apiKey"
)

// APIKeyAuth middleware validates API key and adds it to context
func APIKeyAuth(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract key from header
			key := extractAPIKey(r)
			if key == "" {
				respondUnauthorized(w, "Missing API key. Use 'Authorization: Bearer <key>' or 'X-API-Key: <key>' header.")
				return
			}

			// Hash and lookup
			keyHash := auth.HashAPIKey(key)

			apiKey, err := db.GetApiKeyByHash(keyHash)
			if err != nil {
				respondUnauthorized(w, "Invalid API key")
				return
			}

			if !apiKey.IsActive {
				respondUnauthorized(w, "API key has been revoked")
				return
			}

			// Check expiration
			if apiKey.ExpiresAt.Valid && apiKey.ExpiresAt.Time.Before(time.Now()) {
				respondUnauthorized(w, "API key has expired")
				return
			}

			// Update last used (async)
			go db.UpdateApiKeyLastUsed(apiKey.ID)

			// Add to context for downstream middleware
			ctx := context.WithValue(r.Context(), APIKeyContextKey, apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireSport middleware checks if the API key has access to the requested sport
func RequireSport(sport string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey, ok := r.Context().Value(APIKeyContextKey).(*database.ApiKey)
			if !ok {
				respondUnauthorized(w, "API key not found in context")
				return
			}

			// Check if key has access to this sport
			hasAccess := false
			for _, s := range apiKey.Sports {
				if s == "*" || s == sport {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				respondForbidden(w, "API key does not have access to "+sport+" data")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractAPIKey extracts the API key from the request
func extractAPIKey(r *http.Request) string {
	// Check Authorization header: "Bearer sk_live_..."
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	// Fallback to X-API-Key header
	return r.Header.Get("X-API-Key")
}

// respondUnauthorized sends a 401 response
func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"success":false,"error":{"code":"UNAUTHORIZED","message":"` + message + `"}}`))
}

// respondForbidden sends a 403 response
func respondForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"success":false,"error":{"code":"FORBIDDEN","message":"` + message + `"}}`))
}

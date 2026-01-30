package middleware

import (
	"net/http"
	"sync"

	"github.com/dusanbre/otg-sports-api/internal/database"
	"golang.org/x/time/rate"
)

// RateLimiter implements per-API-key rate limiting
type RateLimiter struct {
	limiters     sync.Map // map[keyHash]*rate.Limiter
	defaultRate  rate.Limit
	defaultBurst int
}

// NewRateLimiter creates a new rate limiter with default requests per minute
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		defaultRate:  rate.Limit(float64(requestsPerMinute) / 60.0), // Convert to per-second
		defaultBurst: requestsPerMinute / 10,                        // Allow small bursts (10% of limit)
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, ok := r.Context().Value(APIKeyContextKey).(*database.ApiKey)
		if !ok {
			// If no API key in context, skip rate limiting (auth middleware should catch this)
			next.ServeHTTP(w, r)
			return
		}

		// Get or create limiter for this API key
		limiter := rl.getLimiter(apiKey)

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"success":false,"error":{"code":"RATE_LIMIT_EXCEEDED","message":"Rate limit exceeded. Please wait before making more requests."}}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getLimiter returns the rate limiter for a specific API key
func (rl *RateLimiter) getLimiter(apiKey *database.ApiKey) *rate.Limiter {
	// Use key hash as identifier
	key := apiKey.KeyHash

	// Try to load existing limiter
	if limiter, ok := rl.limiters.Load(key); ok {
		return limiter.(*rate.Limiter)
	}

	// Calculate rate based on API key's rate limit
	ratePerSecond := rate.Limit(float64(apiKey.RateLimit) / 60.0)
	burst := int(apiKey.RateLimit) / 10
	if burst < 1 {
		burst = 1
	}

	// Create new limiter
	limiter := rate.NewLimiter(ratePerSecond, burst)

	// Store and return (handle race condition)
	actual, _ := rl.limiters.LoadOrStore(key, limiter)
	return actual.(*rate.Limiter)
}

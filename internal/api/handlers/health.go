package handlers

import (
	"net/http"

	"github.com/dusanbre/otg-sports-api/internal/api/middleware"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
//
//	@Summary		Health check
//	@Description	Returns the health status of the API service (public endpoint)
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=map[string]string}
//	@Router			/health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "otg-sport-api",
	})
}

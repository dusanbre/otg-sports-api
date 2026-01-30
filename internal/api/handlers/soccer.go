package handlers

import (
	"net/http"
	"strconv"

	"github.com/dusanbre/otg-sports-api/internal/api/dto"
	"github.com/dusanbre/otg-sports-api/internal/api/middleware"
	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/go-chi/chi/v5"
)

// SoccerHandler handles soccer-related endpoints
type SoccerHandler struct {
	db *database.DB
}

// NewSoccerHandler creates a new soccer handler
func NewSoccerHandler(db *database.DB) *SoccerHandler {
	return &SoccerHandler{db: db}
}

// GetMatches returns a list of soccer matches
func (h *SoccerHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params := parseQueryParams(r)

	// Fetch matches from database
	matches, total, err := h.db.GetSoccerMatchesFiltered(params)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch matches")
		return
	}

	// Convert to DTOs
	response := make([]dto.SoccerMatchResponse, len(matches))
	for i, m := range matches {
		response[i] = dto.SoccerMatchFromModel(&m)
	}

	middleware.RespondJSONWithMeta(w, http.StatusOK, response, &middleware.MetaInfo{
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
}

// GetMatch returns a single soccer match by ID
func (h *SoccerHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid match ID")
		return
	}

	match, err := h.db.GetMatchByID(id)
	if err != nil {
		middleware.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Match not found")
		return
	}

	response := dto.SoccerMatchFromModel(match)
	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetLiveMatches returns currently live soccer matches
func (h *SoccerHandler) GetLiveMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.db.GetLiveSoccerMatches()
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch live matches")
		return
	}

	response := make([]dto.SoccerMatchResponse, len(matches))
	for i, m := range matches {
		response[i] = dto.SoccerMatchFromModel(&m)
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetLeagues returns a list of available soccer leagues
func (h *SoccerHandler) GetLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := h.db.GetSoccerLeagues()
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch leagues")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, leagues)
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/dusanbre/otg-sports-api/internal/api/dto"
	"github.com/dusanbre/otg-sports-api/internal/api/middleware"
	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/go-chi/chi/v5"
)

// BasketballHandler handles basketball-related endpoints
type BasketballHandler struct {
	db *database.DB
}

// NewBasketballHandler creates a new basketball handler
func NewBasketballHandler(db *database.DB) *BasketballHandler {
	return &BasketballHandler{db: db}
}

// GetMatches returns a list of basketball matches
func (h *BasketballHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params := parseQueryParams(r)

	// Fetch matches from database
	matches, total, err := h.db.GetBasketballMatchesFiltered(params)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch matches")
		return
	}

	// Convert to DTOs
	response := make([]dto.BasketballMatchResponse, len(matches))
	for i, m := range matches {
		response[i] = dto.BasketballMatchFromModel(&m)
	}

	middleware.RespondJSONWithMeta(w, http.StatusOK, response, &middleware.MetaInfo{
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
}

// GetMatch returns a single basketball match by ID
func (h *BasketballHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid match ID")
		return
	}

	match, err := h.db.GetBasketballMatchByID(id)
	if err != nil {
		middleware.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Match not found")
		return
	}

	response := dto.BasketballMatchFromModel(match)
	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetLiveMatches returns currently live basketball matches
func (h *BasketballHandler) GetLiveMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.db.GetLiveBasketballMatches()
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch live matches")
		return
	}

	response := make([]dto.BasketballMatchResponse, len(matches))
	for i, m := range matches {
		response[i] = dto.BasketballMatchFromModel(&m)
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetLeagues returns a list of available basketball leagues
func (h *BasketballHandler) GetLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := h.db.GetBasketballLeagues()
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch leagues")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, leagues)
}

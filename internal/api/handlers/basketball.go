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

// GetMatches godoc
//
//	@Summary		List basketball matches
//	@Description	Returns a paginated list of basketball matches with optional filtering
//	@Tags			basketball
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int		false	"Maximum results (1-100)"	default(50)
//	@Param			offset		query		int		false	"Results to skip"			default(0)
//	@Param			date		query		string	false	"Filter by date (YYYY-MM-DD)"
//	@Param			status		query		string	false	"Filter by status (FT, Q1, Q2, Q3, Q4, OT, HT)"
//	@Param			league_id	query		int		false	"Filter by league ID"
//	@Success		200			{object}	middleware.Response{data=[]dto.BasketballMatchResponse,meta=middleware.MetaInfo}
//	@Failure		401			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		403			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		429			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		500			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/basketball/matches [get]
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

// GetMatch godoc
//
//	@Summary		Get basketball match by ID
//	@Description	Returns a single basketball match by its match ID
//	@Tags			basketball
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Match ID"
//	@Success		200	{object}	middleware.Response{data=dto.BasketballMatchResponse}
//	@Failure		400	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		401	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		404	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/basketball/matches/{id} [get]
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

// GetLiveMatches godoc
//
//	@Summary		Get live basketball matches
//	@Description	Returns all currently live basketball matches
//	@Tags			basketball
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=[]dto.BasketballMatchResponse}
//	@Failure		401	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		403	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		500	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/basketball/matches/live [get]
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

// GetLeagues godoc
//
//	@Summary		Get basketball leagues
//	@Description	Returns a list of all available basketball leagues
//	@Tags			basketball
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=[]database.LeagueInfo}
//	@Failure		401	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		403	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		500	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/basketball/leagues [get]
func (h *BasketballHandler) GetLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := h.db.GetBasketballLeagues()
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch leagues")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, leagues)
}

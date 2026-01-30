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

// GetMatches godoc
//
//	@Summary		List soccer matches
//	@Description	Returns a paginated list of soccer matches with optional filtering
//	@Tags			soccer
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int		false	"Maximum results (1-100)"	default(50)
//	@Param			offset		query		int		false	"Results to skip"			default(0)
//	@Param			date		query		string	false	"Filter by date (YYYY-MM-DD)"
//	@Param			status		query		string	false	"Filter by status (FT, 1H, HT, 2H, NS)"
//	@Param			league_id	query		int		false	"Filter by league ID"
//	@Success		200			{object}	middleware.Response{data=[]dto.SoccerMatchResponse,meta=middleware.MetaInfo}
//	@Failure		401			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		403			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		429			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		500			{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/soccer/matches [get]
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

// GetMatch godoc
//
//	@Summary		Get soccer match by ID
//	@Description	Returns a single soccer match by its match ID
//	@Tags			soccer
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Match ID"
//	@Success		200	{object}	middleware.Response{data=dto.SoccerMatchResponse}
//	@Failure		400	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		401	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		404	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/soccer/matches/{id} [get]
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

// GetLiveMatches godoc
//
//	@Summary		Get live soccer matches
//	@Description	Returns all currently live soccer matches
//	@Tags			soccer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=[]dto.SoccerMatchResponse}
//	@Failure		401	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		403	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		500	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/soccer/matches/live [get]
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

// GetLeagues godoc
//
//	@Summary		Get soccer leagues
//	@Description	Returns a list of all available soccer leagues
//	@Tags			soccer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=[]database.LeagueInfo}
//	@Failure		401	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		403	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Failure		500	{object}	middleware.Response{error=middleware.ErrorInfo}
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
//	@Router			/soccer/leagues [get]
func (h *SoccerHandler) GetLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := h.db.GetSoccerLeagues()
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch leagues")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, leagues)
}

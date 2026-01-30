package handlers

import (
	"net/http"
	"strconv"

	"github.com/dusanbre/otg-sports-api/internal/database"
)

// parseQueryParams extracts common query parameters from the request
func parseQueryParams(r *http.Request) database.QueryParams {
	params := database.QueryParams{
		Limit:  50, // Default limit
		Offset: 0,
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			params.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}

	// Parse date filter
	params.Date = r.URL.Query().Get("date")

	// Parse status filter
	params.Status = r.URL.Query().Get("status")

	// Parse league_id filter
	if leagueIDStr := r.URL.Query().Get("league_id"); leagueIDStr != "" {
		if leagueID, err := strconv.ParseInt(leagueIDStr, 10, 64); err == nil {
			params.LeagueID = &leagueID
		}
	}

	return params
}

package dto

import "github.com/dusanbre/otg-sports-api/internal/database"

// TeamInfo represents a team in the API response
type TeamInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Score *int   `json:"score,omitempty"`
}

// SoccerMatchResponse is the API response for a soccer match
type SoccerMatchResponse struct {
	ID            int64    `json:"id"`
	MatchID       int64    `json:"match_id"`
	Sport         string   `json:"sport"`
	LeagueID      int64    `json:"league_id"`
	LeagueGID     int64    `json:"league_gid"`
	LeagueName    string   `json:"league_name"`
	Status        string   `json:"status"`
	StartDate     string   `json:"start_date"`
	StartTime     string   `json:"start_time"`
	HomeTeam      TeamInfo `json:"home_team"`
	AwayTeam      TeamInfo `json:"away_team"`
	HalfTimeScore string   `json:"half_time_score,omitempty"`
	FullTimeScore string   `json:"full_time_score,omitempty"`
}

// SoccerMatchFromModel converts a database model to API response
func SoccerMatchFromModel(m *database.SoccerMatch) SoccerMatchResponse {
	response := SoccerMatchResponse{
		ID:         m.ID,
		Sport:      "soccer",
		LeagueName: m.LeagueName.String,
		Status:     m.MatchStatus.String,
	}

	if m.MatchID.Valid {
		response.MatchID = m.MatchID.Int64
	}
	if m.LeagueID.Valid {
		response.LeagueID = m.LeagueID.Int64
	}
	if m.LeagueGID.Valid {
		response.LeagueGID = m.LeagueGID.Int64
	}
	if m.MatchStartDate.Valid {
		response.StartDate = m.MatchStartDate.Time.Format("2006-01-02")
	}
	if m.MatchStartTime.Valid {
		response.StartTime = m.MatchStartTime.String
	}

	// Home team
	response.HomeTeam = TeamInfo{
		Name: m.HTeamName.String,
	}
	if m.HTeamID.Valid {
		response.HomeTeam.ID = m.HTeamID.Int64
	}
	if m.HTeamGoals.Valid {
		goals := int(m.HTeamGoals.Int32)
		response.HomeTeam.Score = &goals
	}

	// Away team
	response.AwayTeam = TeamInfo{
		Name: m.ATeamName.String,
	}
	if m.ATeamID.Valid {
		response.AwayTeam.ID = m.ATeamID.Int64
	}
	if m.ATeamGoals.Valid {
		goals := int(m.ATeamGoals.Int32)
		response.AwayTeam.Score = &goals
	}

	// Scores
	if m.HTScore.Valid {
		response.HalfTimeScore = m.HTScore.String
	}
	if m.FTScore.Valid {
		response.FullTimeScore = m.FTScore.String
	}

	return response
}

// LeagueInfo represents a league in the API response
type LeagueInfo struct {
	ID   int64  `json:"id"`
	GID  int64  `json:"gid"`
	Name string `json:"name"`
}

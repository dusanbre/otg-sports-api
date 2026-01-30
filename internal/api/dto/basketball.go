package dto

import "github.com/dusanbre/otg-sports-api/internal/database"

// QuarterScores represents quarter-by-quarter scores for basketball
type QuarterScores struct {
	Q1 *ScorePair `json:"q1,omitempty"`
	Q2 *ScorePair `json:"q2,omitempty"`
	Q3 *ScorePair `json:"q3,omitempty"`
	Q4 *ScorePair `json:"q4,omitempty"`
	OT *ScorePair `json:"ot,omitempty"`
}

// ScorePair represents home/away scores
type ScorePair struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// BasketballMatchResponse is the API response for a basketball match
type BasketballMatchResponse struct {
	ID            int64          `json:"id"`
	MatchID       int64          `json:"match_id"`
	Sport         string         `json:"sport"`
	LeagueID      int64          `json:"league_id"`
	LeagueGID     int64          `json:"league_gid"`
	LeagueName    string         `json:"league_name"`
	FileGroup     string         `json:"file_group,omitempty"`
	Status        string         `json:"status"`
	StartDate     string         `json:"start_date"`
	StartTime     string         `json:"start_time"`
	Timer         string         `json:"timer,omitempty"`
	HomeTeam      TeamInfo       `json:"home_team"`
	AwayTeam      TeamInfo       `json:"away_team"`
	QuarterScores *QuarterScores `json:"quarter_scores,omitempty"`
}

// BasketballMatchFromModel converts a database model to API response
func BasketballMatchFromModel(m *database.BasketballMatch) BasketballMatchResponse {
	response := BasketballMatchResponse{
		ID:         m.ID,
		Sport:      "basketball",
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
	if m.FileGroup.Valid {
		response.FileGroup = m.FileGroup.String
	}
	if m.MatchDate.Valid {
		response.StartDate = m.MatchDate.Time.Format("2006-01-02")
	}
	if m.MatchTime.Valid {
		response.StartTime = m.MatchTime.String
	}
	if m.Timer.Valid {
		response.Timer = m.Timer.String
	}

	// Home team
	response.HomeTeam = TeamInfo{
		Name: m.HTeamName.String,
	}
	if m.HTeamID.Valid {
		response.HomeTeam.ID = m.HTeamID.Int64
	}
	if m.HTeamScore.Valid {
		score := int(m.HTeamScore.Int32)
		response.HomeTeam.Score = &score
	}

	// Away team
	response.AwayTeam = TeamInfo{
		Name: m.ATeamName.String,
	}
	if m.ATeamID.Valid {
		response.AwayTeam.ID = m.ATeamID.Int64
	}
	if m.ATeamScore.Valid {
		score := int(m.ATeamScore.Int32)
		response.AwayTeam.Score = &score
	}

	// Quarter scores
	hasQuarterData := m.HTeamQ1.Valid || m.HTeamQ2.Valid || m.HTeamQ3.Valid || m.HTeamQ4.Valid || m.HTeamOt.Valid
	if hasQuarterData {
		response.QuarterScores = &QuarterScores{}

		if m.HTeamQ1.Valid && m.ATeamQ1.Valid {
			response.QuarterScores.Q1 = &ScorePair{
				Home: int(m.HTeamQ1.Int32),
				Away: int(m.ATeamQ1.Int32),
			}
		}
		if m.HTeamQ2.Valid && m.ATeamQ2.Valid {
			response.QuarterScores.Q2 = &ScorePair{
				Home: int(m.HTeamQ2.Int32),
				Away: int(m.ATeamQ2.Int32),
			}
		}
		if m.HTeamQ3.Valid && m.ATeamQ3.Valid {
			response.QuarterScores.Q3 = &ScorePair{
				Home: int(m.HTeamQ3.Int32),
				Away: int(m.ATeamQ3.Int32),
			}
		}
		if m.HTeamQ4.Valid && m.ATeamQ4.Valid {
			response.QuarterScores.Q4 = &ScorePair{
				Home: int(m.HTeamQ4.Int32),
				Away: int(m.ATeamQ4.Int32),
			}
		}
		if m.HTeamOt.Valid && m.ATeamOt.Valid {
			response.QuarterScores.OT = &ScorePair{
				Home: int(m.HTeamOt.Int32),
				Away: int(m.ATeamOt.Int32),
			}
		}
	}

	return response
}

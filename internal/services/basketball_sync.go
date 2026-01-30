package services

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/dusanbre/otg-sports-api/internal/goalserve"
)

// BasketballSyncService handles syncing basketball matches from Goalserve to database
type BasketballSyncService struct {
	db              *database.DB
	goalserveClient *goalserve.Client
}

// NewBasketballSyncService creates a new basketball sync service
func NewBasketballSyncService(db *database.DB) *BasketballSyncService {
	return &BasketballSyncService{
		db:              db,
		goalserveClient: goalserve.NewClient(),
	}
}

// SyncMatches fetches basketball matches from Goalserve and syncs them to the database
func (s *BasketballSyncService) SyncMatches() error {
	log.Println("Starting basketball match sync...")

	// Fetch today's matches
	basketballData, err := s.goalserveClient.FetchBasketballTodayMatches()
	if err != nil {
		return fmt.Errorf("failed to fetch today's basketball matches from Goalserve: %w", err)
	}

	matchesInserted := 0
	matchesUpdated := 0

	// Process today's matches
	for _, category := range basketballData.Categories {
		for _, match := range category.Match.Matches {
			inserted, err := s.upsertBasketballMatch(category, match)
			if err != nil {
				log.Printf("Failed to upsert basketball match %s: %v", match.ID, err)
				continue
			}
			if inserted {
				matchesInserted++
			} else {
				matchesUpdated++
			}
		}
	}

	// Fetch future matches (next 7 days)
	futureData, err := s.goalserveClient.FetchBasketballMatchesFuture7Days()
	if err != nil {
		log.Printf("Warning: failed to fetch future basketball matches: %v", err)
	} else {
		// Process future matches
		for _, category := range futureData.Categories {
			for _, match := range category.Match.Matches {
				inserted, err := s.upsertBasketballMatch(category, match)
				if err != nil {
					log.Printf("Failed to upsert future basketball match %s: %v", match.ID, err)
					continue
				}
				if inserted {
					matchesInserted++
				} else {
					matchesUpdated++
				}
			}
		}
	}

	log.Printf("Basketball match sync completed: %d inserted, %d updated", matchesInserted, matchesUpdated)
	return nil
}

// upsertBasketballMatch inserts or updates a basketball match in the database
func (s *BasketballSyncService) upsertBasketballMatch(category goalserve.GoalServeBasketballCategory, match goalserve.GoalServeBasketballMatch) (bool, error) {
	// Parse match ID
	matchID, err := strconv.ParseInt(match.ID, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid match ID: %w", err)
	}

	// Parse league IDs
	leagueID, _ := strconv.ParseInt(category.ID, 10, 64)
	leagueGid, _ := strconv.ParseInt(category.Gid, 10, 64)

	// Parse team IDs
	hTeamID, _ := strconv.ParseInt(match.LocalTeam.ID, 10, 64)
	aTeamID, _ := strconv.ParseInt(match.AwayTeam.ID, 10, 64)

	// Parse date and time
	dateStr := match.Date
	if dateStr == "" || match.Time == "" {
		return false, fmt.Errorf("missing date or time data: date='%s', time='%s'", dateStr, match.Time)
	}

	// Parse date (format: "29.01.2026")
	matchDate, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return false, fmt.Errorf("invalid date format: %s", dateStr)
	}

	// Parse time (format: "23:30")
	timeOnly, err := time.Parse("15:04", match.Time)
	if err != nil {
		return false, fmt.Errorf("invalid time format: %s", match.Time)
	}

	// Create time value for TIME column
	year, month, day := time.Now().Date()
	matchTime := time.Date(year, month, day, timeOnly.Hour(), timeOnly.Minute(), 0, 0, time.UTC)

	// Parse scores - helper function for nullable int32
	parseScore := func(s string) sql.NullInt32 {
		if s == "" {
			return sql.NullInt32{Valid: false}
		}
		if score, err := strconv.Atoi(s); err == nil {
			return sql.NullInt32{Int32: int32(score), Valid: true}
		}
		return sql.NullInt32{Valid: false}
	}

	// Parse team scores
	hTeamScore := parseScore(match.LocalTeam.TotalScore)
	hTeamQ1 := parseScore(match.LocalTeam.Q1)
	hTeamQ2 := parseScore(match.LocalTeam.Q2)
	hTeamQ3 := parseScore(match.LocalTeam.Q3)
	hTeamQ4 := parseScore(match.LocalTeam.Q4)
	hTeamOt := parseScore(match.LocalTeam.Ot)

	aTeamScore := parseScore(match.AwayTeam.TotalScore)
	aTeamQ1 := parseScore(match.AwayTeam.Q1)
	aTeamQ2 := parseScore(match.AwayTeam.Q2)
	aTeamQ3 := parseScore(match.AwayTeam.Q3)
	aTeamQ4 := parseScore(match.AwayTeam.Q4)
	aTeamOt := parseScore(match.AwayTeam.Ot)

	// Timer
	var timer sql.NullString
	if match.Timer != "" {
		timer = sql.NullString{String: match.Timer, Valid: true}
	}

	// Check if match exists
	var existingID int64
	checkQuery := s.db.Builder.
		Select("id").
		From("basketball_matches").
		Where("match_id = ?", matchID)

	checkSQL, checkArgs, _ := checkQuery.ToSql()
	err = s.db.Conn.QueryRow(checkSQL, checkArgs...).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new match
		insertQuery := s.db.Builder.
			Insert("basketball_matches").
			Columns(
				"match_id", "league_gid", "league_id", "league_name", "file_group",
				"match_status", "match_date", "match_time", "timer",
				"h_team_id", "h_team_name", "h_team_score",
				"h_team_q1", "h_team_q2", "h_team_q3", "h_team_q4", "h_team_ot",
				"a_team_id", "a_team_name", "a_team_score",
				"a_team_q1", "a_team_q2", "a_team_q3", "a_team_q4", "a_team_ot",
			).
			Values(
				matchID, leagueGid, leagueID, category.Name, category.FileGroup,
				match.Status, matchDate, matchTime, timer,
				hTeamID, match.LocalTeam.Name, hTeamScore,
				hTeamQ1, hTeamQ2, hTeamQ3, hTeamQ4, hTeamOt,
				aTeamID, match.AwayTeam.Name, aTeamScore,
				aTeamQ1, aTeamQ2, aTeamQ3, aTeamQ4, aTeamOt,
			)

		insertSQL, insertArgs, err := insertQuery.ToSql()
		if err != nil {
			return false, fmt.Errorf("failed to build insert query: %w", err)
		}

		_, err = s.db.Conn.Exec(insertSQL, insertArgs...)
		if err != nil {
			return false, fmt.Errorf("failed to insert basketball match: %w", err)
		}

		log.Printf("Inserted basketball match: %s vs %s", match.LocalTeam.Name, match.AwayTeam.Name)
		return true, nil
	} else if err == nil {
		// Update existing match
		updateQuery := s.db.Builder.
			Update("basketball_matches").
			Set("match_status", match.Status).
			Set("timer", timer).
			Set("h_team_score", hTeamScore).
			Set("h_team_q1", hTeamQ1).
			Set("h_team_q2", hTeamQ2).
			Set("h_team_q3", hTeamQ3).
			Set("h_team_q4", hTeamQ4).
			Set("h_team_ot", hTeamOt).
			Set("a_team_score", aTeamScore).
			Set("a_team_q1", aTeamQ1).
			Set("a_team_q2", aTeamQ2).
			Set("a_team_q3", aTeamQ3).
			Set("a_team_q4", aTeamQ4).
			Set("a_team_ot", aTeamOt).
			Set("updated_at", time.Now()).
			Where("match_id = ?", matchID)

		updateSQL, updateArgs, err := updateQuery.ToSql()
		if err != nil {
			return false, fmt.Errorf("failed to build update query: %w", err)
		}

		_, err = s.db.Conn.Exec(updateSQL, updateArgs...)
		if err != nil {
			return false, fmt.Errorf("failed to update basketball match: %w", err)
		}

		log.Printf("Updated basketball match: %s vs %s", match.LocalTeam.Name, match.AwayTeam.Name)
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check if basketball match exists: %w", err)
	}
}

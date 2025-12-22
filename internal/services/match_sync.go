package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/dusanbre/otg-sports-api/internal/goalserve"
)

// MatchSyncService handles syncing matches from Goalserve to database
type MatchSyncService struct {
	db              *database.DB
	goalserveClient *goalserve.Client
}

// NewMatchSyncService creates a new match sync service
func NewMatchSyncService(db *database.DB) *MatchSyncService {
	return &MatchSyncService{
		db:              db,
		goalserveClient: goalserve.NewClient(),
	}
}

// SyncMatches fetches matches from Goalserve and syncs them to the database
func (s *MatchSyncService) SyncMatches() error {
	log.Println("Starting match sync...")

	// Fetch today's matches
	soccerData, err := s.goalserveClient.FetchTodayMatches()
	if err != nil {
		return fmt.Errorf("failed to fetch today's matches from Goalserve: %w", err)
	}

	matchesInserted := 0
	matchesUpdated := 0

	// Process today's matches
	for _, category := range soccerData.Categories {
		for _, match := range category.Matches.Match {
			inserted, err := s.upsertMatch(category, match)
			if err != nil {
				log.Printf("Failed to upsert match %s: %v", match.ID, err)
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
	futureData, err := s.goalserveClient.FetchMatchesFuture7Days()
	if err != nil {
		log.Printf("Warning: failed to fetch future matches: %v", err)
	} else {
		// Process future matches
		for _, category := range futureData.Categories {
			for _, match := range category.Matches.Match {
				inserted, err := s.upsertMatch(category, match)
				if err != nil {
					log.Printf("Failed to upsert future match %s: %v", match.ID, err)
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

	log.Printf("Match sync completed: %d inserted, %d updated", matchesInserted, matchesUpdated)
	return nil
}

// upsertMatch inserts or updates a match in the database
func (s *MatchSyncService) upsertMatch(category goalserve.GoalServeScoreCategory, match goalserve.GoalServeMatch) (bool, error) {
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
	aTeamID, _ := strconv.ParseInt(match.VisitorTeam.ID, 10, 64)

	// Parse date and time - use FormattedDate if available, otherwise fallback to Date
	dateStr := match.FormattedDate
	if dateStr == "" {
		dateStr = match.Date
	}

	// Parse combined date and time
	var matchDate, matchTime time.Time
	if dateStr == "" || match.Time == "" {
		return false, fmt.Errorf("missing date or time data: date='%s', time='%s'", dateStr, match.Time)
	}

	// Try to parse combined date and time
	dateTimeStr := dateStr + " " + match.Time
	combinedDateTime, err := time.Parse("02.01.2006 15:04", dateTimeStr)
	if err != nil {
		// Try alternative date format (e.g., "Dec 22 2025 15:04")
		combinedDateTime, err = time.Parse("Jan 2 2006 15:04", dateStr+" "+fmt.Sprintf("%d", time.Now().Year())+" "+match.Time)
	}

	if err == nil {
		matchDate = combinedDateTime.Truncate(24 * time.Hour) // Date only (start of day)
		// Extract only time portion for TIME column
		year, month, day := time.Now().Date()
		matchTime = time.Date(year, month, day, combinedDateTime.Hour(), combinedDateTime.Minute(), combinedDateTime.Second(), 0, time.UTC)
	} else {
		// Try separate parsing as fallback
		matchDate, err = time.Parse("02.01.2006", dateStr)
		if err != nil {
			return false, fmt.Errorf("invalid date format: %s", dateStr)
		}

		// For time, extract only time portion for TIME column
		timeOnly, err := time.Parse("15:04", match.Time)
		if err != nil {
			return false, fmt.Errorf("invalid time format: %s", match.Time)
		}

		// Use consistent date reference with parsed time
		year, month, day := time.Now().Date()
		matchTime = time.Date(year, month, day, timeOnly.Hour(), timeOnly.Minute(), 0, 0, time.UTC)
	}

	// Parse goals
	var hGoals, aGoals sql.NullInt32
	if match.LocalTeam.Goals != "" {
		if goals, err := strconv.Atoi(match.LocalTeam.Goals); err == nil {
			hGoals = sql.NullInt32{Int32: int32(goals), Valid: true}
		}
	}

	if match.VisitorTeam.Goals != "" {
		if goals, err := strconv.Atoi(match.VisitorTeam.Goals); err == nil {
			aGoals = sql.NullInt32{Int32: int32(goals), Valid: true}
		}
	}

	// Set scores
	var htScore, ftScore sql.NullString
	if match.HTScore.Score != "" {
		htScore = sql.NullString{String: match.HTScore.Score, Valid: true}
	}

	if match.FTScore.Score != "" {
		ftScore = sql.NullString{String: match.FTScore.Score, Valid: true}
	}

	// Convert events to JSON - handle events which can be null or an object with event array
	eventsJSON := []byte("[]")
	if match.Events != nil {
		// Try to extract events if it's an object with event array
		eventsBytes, err := json.Marshal(match.Events)
		if err == nil {
			var eventsWrapper struct {
				Event []goalserve.GoalServeEvent `json:"event"`
			}
			if err := json.Unmarshal(eventsBytes, &eventsWrapper); err == nil && len(eventsWrapper.Event) > 0 {
				eventsJSON, _ = json.Marshal(eventsWrapper.Event)
			}
		}
	}

	// Check if match exists
	var existingID int64
	checkQuery := s.db.Builder.
		Select("id").
		From("soccer_matches").
		Where("match_id = ?", matchID)

	checkSQL, checkArgs, _ := checkQuery.ToSql()
	err = s.db.Conn.QueryRow(checkSQL, checkArgs...).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new match
		insertQuery := s.db.Builder.
			Insert("soccer_matches").
			Columns(
				"match_id", "league_gid", "league_id", "league_name",
				"match_status", "match_start_date", "match_start_time",
				"h_team_id", "a_team_id", "h_team_name", "a_team_name",
				"h_team_goals", "a_team_goals", "ht_score", "ft_score", "events",
			).
			Values(
				matchID, leagueGid, leagueID, category.Name,
				match.Status, matchDate, matchTime,
				hTeamID, aTeamID, match.LocalTeam.Name, match.VisitorTeam.Name,
				hGoals, aGoals, htScore, ftScore, string(eventsJSON),
			)

		insertSQL, insertArgs, err := insertQuery.ToSql()
		if err != nil {
			return false, fmt.Errorf("failed to build insert query: %w", err)
		}

		_, err = s.db.Conn.Exec(insertSQL, insertArgs...)
		if err != nil {
			return false, fmt.Errorf("failed to insert match: %w", err)
		}

		log.Printf("Inserted match: %s vs %s", match.LocalTeam.Name, match.VisitorTeam.Name)
		return true, nil
	} else if err == nil {
		// Update existing match
		updateQuery := s.db.Builder.
			Update("soccer_matches").
			Set("match_status", match.Status).
			Set("h_team_goals", hGoals).
			Set("a_team_goals", aGoals).
			Set("ht_score", htScore).
			Set("ft_score", ftScore).
			Set("events", string(eventsJSON)).
			Set("updated_at", time.Now()).
			Where("match_id = ?", matchID)

		updateSQL, updateArgs, err := updateQuery.ToSql()
		if err != nil {
			return false, fmt.Errorf("failed to build update query: %w", err)
		}

		_, err = s.db.Conn.Exec(updateSQL, updateArgs...)
		if err != nil {
			return false, fmt.Errorf("failed to update match: %w", err)
		}

		log.Printf("Updated match: %s vs %s", match.LocalTeam.Name, match.VisitorTeam.Name)
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check if match exists: %w", err)
	}
}

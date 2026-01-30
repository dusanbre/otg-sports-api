package database

import (
	"encoding/json"
	"fmt"
	"time"
)

// Example: Query soccer matches
func (db *DB) GetSoccerMatches() ([]SoccerMatch, error) {
	query := db.Builder.
		Select("*").
		From("soccer_matches").
		OrderBy("match_start_date DESC").
		Limit(10)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var matches []SoccerMatch
	for rows.Next() {
		var m SoccerMatch
		err := rows.Scan(
			&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName,
			&m.MatchStatus, &m.MatchStartDate, &m.MatchStartTime,
			&m.HTeamID, &m.ATeamID, &m.HTeamName, &m.ATeamName,
			&m.HTeamGoals, &m.ATeamGoals, &m.HTScore, &m.FTScore,
			&m.Events, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		matches = append(matches, m)
	}

	return matches, nil
}

// Example: Get match by ID
func (db *DB) GetMatchByID(matchID int64) (*SoccerMatch, error) {
	query := db.Builder.
		Select("*").
		From("soccer_matches").
		Where("match_id = ?", matchID)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var m SoccerMatch
	err = db.Conn.QueryRow(sql, args...).Scan(
		&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName,
		&m.MatchStatus, &m.MatchStartDate, &m.MatchStartTime,
		&m.HTeamID, &m.ATeamID, &m.HTeamName, &m.ATeamName,
		&m.HTeamGoals, &m.ATeamGoals, &m.HTScore, &m.FTScore,
		&m.Events, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query match: %w", err)
	}

	return &m, nil
}

// Example: Insert a soccer match
func (db *DB) InsertMatch(m *SoccerMatch) error {
	query := db.Builder.
		Insert("soccer_matches").
		Columns(
			"match_id", "league_gid", "league_id", "league_name",
			"match_status", "match_start_date", "match_start_time",
			"h_team_id", "a_team_id", "h_team_name", "a_team_name",
			"h_team_goals", "a_team_goals", "ht_score", "ft_score", "events",
		).
		Values(
			m.MatchID, m.LeagueGID, m.LeagueID, m.LeagueName,
			m.MatchStatus, m.MatchStartDate, m.MatchStartTime,
			m.HTeamID, m.ATeamID, m.HTeamName, m.ATeamName,
			m.HTeamGoals, m.ATeamGoals, m.HTScore, m.FTScore, m.Events,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	err = db.Conn.QueryRow(sql, args...).Scan(&m.ID)
	if err != nil {
		return fmt.Errorf("failed to insert match: %w", err)
	}

	return nil
}

// Example: Update match score
func (db *DB) UpdateMatchScore(matchID int64, hGoals, aGoals int, htScore, ftScore string) error {
	query := db.Builder.
		Update("soccer_matches").
		Set("h_team_goals", hGoals).
		Set("a_team_goals", aGoals).
		Set("ht_score", htScore).
		Set("ft_score", ftScore).
		Where("match_id = ?", matchID)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = db.Conn.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update match: %w", err)
	}

	return nil
}

// GetBasketballMatches returns recent basketball matches
func (db *DB) GetBasketballMatches() ([]BasketballMatch, error) {
	query := db.Builder.
		Select("*").
		From("basketball_matches").
		OrderBy("match_date DESC").
		Limit(10)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var matches []BasketballMatch
	for rows.Next() {
		var m BasketballMatch
		err := rows.Scan(
			&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName, &m.FileGroup,
			&m.MatchStatus, &m.MatchDate, &m.MatchTime, &m.Timer,
			&m.HTeamID, &m.HTeamName, &m.HTeamScore,
			&m.HTeamQ1, &m.HTeamQ2, &m.HTeamQ3, &m.HTeamQ4, &m.HTeamOt,
			&m.ATeamID, &m.ATeamName, &m.ATeamScore,
			&m.ATeamQ1, &m.ATeamQ2, &m.ATeamQ3, &m.ATeamQ4, &m.ATeamOt,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		matches = append(matches, m)
	}

	return matches, nil
}

// GetBasketballMatchByID returns a basketball match by its match ID
func (db *DB) GetBasketballMatchByID(matchID int64) (*BasketballMatch, error) {
	query := db.Builder.
		Select("*").
		From("basketball_matches").
		Where("match_id = ?", matchID)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var m BasketballMatch
	err = db.Conn.QueryRow(sql, args...).Scan(
		&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName, &m.FileGroup,
		&m.MatchStatus, &m.MatchDate, &m.MatchTime, &m.Timer,
		&m.HTeamID, &m.HTeamName, &m.HTeamScore,
		&m.HTeamQ1, &m.HTeamQ2, &m.HTeamQ3, &m.HTeamQ4, &m.HTeamOt,
		&m.ATeamID, &m.ATeamName, &m.ATeamScore,
		&m.ATeamQ1, &m.ATeamQ2, &m.ATeamQ3, &m.ATeamQ4, &m.ATeamOt,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query basketball match: %w", err)
	}

	return &m, nil
}

// ============================================================================
// API Key Queries
// ============================================================================

// CreateApiKey creates a new API key record
func (db *DB) CreateApiKey(keyHash, keyPrefix, name string, sports []string, rateLimit int) (*ApiKey, error) {
	sportsJSON, err := json.Marshal(sports)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sports: %w", err)
	}

	query := db.Builder.
		Insert("api_keys").
		Columns("key_hash", "key_prefix", "name", "sports", "rate_limit").
		Values(keyHash, keyPrefix, name, string(sportsJSON), rateLimit).
		Suffix("RETURNING id, key_hash, key_prefix, name, sports, rate_limit, is_active, created_at, last_used_at, expires_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var apiKey ApiKey
	var sportsStr string
	err = db.Conn.QueryRow(sql, args...).Scan(
		&apiKey.ID, &apiKey.KeyHash, &apiKey.KeyPrefix, &apiKey.Name,
		&sportsStr, &apiKey.RateLimit, &apiKey.IsActive,
		&apiKey.CreatedAt, &apiKey.LastUsedAt, &apiKey.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert api key: %w", err)
	}

	// Unmarshal sports JSON
	if err := json.Unmarshal([]byte(sportsStr), &apiKey.Sports); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sports: %w", err)
	}

	return &apiKey, nil
}

// GetApiKeyByHash retrieves an API key by its hash
func (db *DB) GetApiKeyByHash(keyHash string) (*ApiKey, error) {
	query := db.Builder.
		Select("id", "key_hash", "key_prefix", "name", "sports", "rate_limit", "is_active", "created_at", "last_used_at", "expires_at").
		From("api_keys").
		Where("key_hash = ?", keyHash)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var apiKey ApiKey
	var sportsStr string
	err = db.Conn.QueryRow(sql, args...).Scan(
		&apiKey.ID, &apiKey.KeyHash, &apiKey.KeyPrefix, &apiKey.Name,
		&sportsStr, &apiKey.RateLimit, &apiKey.IsActive,
		&apiKey.CreatedAt, &apiKey.LastUsedAt, &apiKey.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query api key: %w", err)
	}

	// Unmarshal sports JSON
	if err := json.Unmarshal([]byte(sportsStr), &apiKey.Sports); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sports: %w", err)
	}

	return &apiKey, nil
}

// GetAllApiKeys retrieves all API keys
func (db *DB) GetAllApiKeys() ([]ApiKey, error) {
	query := db.Builder.
		Select("id", "key_hash", "key_prefix", "name", "sports", "rate_limit", "is_active", "created_at", "last_used_at", "expires_at").
		From("api_keys").
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var keys []ApiKey
	for rows.Next() {
		var apiKey ApiKey
		var sportsStr string
		err := rows.Scan(
			&apiKey.ID, &apiKey.KeyHash, &apiKey.KeyPrefix, &apiKey.Name,
			&sportsStr, &apiKey.RateLimit, &apiKey.IsActive,
			&apiKey.CreatedAt, &apiKey.LastUsedAt, &apiKey.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Unmarshal sports JSON
		if err := json.Unmarshal([]byte(sportsStr), &apiKey.Sports); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sports: %w", err)
		}

		keys = append(keys, apiKey)
	}

	return keys, nil
}

// RevokeApiKey deactivates an API key by ID
func (db *DB) RevokeApiKey(id int64) error {
	query := db.Builder.
		Update("api_keys").
		Set("is_active", false).
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := db.Conn.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke api key: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("api key not found")
	}

	return nil
}

// UpdateApiKeyLastUsed updates the last_used_at timestamp
func (db *DB) UpdateApiKeyLastUsed(id int64) error {
	query := db.Builder.
		Update("api_keys").
		Set("last_used_at", time.Now()).
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = db.Conn.Exec(sql, args...)
	return err
}

// ============================================================================
// Soccer Filtered Queries (for API)
// ============================================================================

// GetSoccerMatchesFiltered returns soccer matches with filtering and pagination
func (db *DB) GetSoccerMatchesFiltered(params QueryParams) ([]SoccerMatch, int, error) {
	// Build base query
	baseQuery := db.Builder.
		Select("*").
		From("soccer_matches")

	countQuery := db.Builder.
		Select("COUNT(*)").
		From("soccer_matches")

	// Apply filters
	if params.Date != "" {
		baseQuery = baseQuery.Where("match_start_date = ?", params.Date)
		countQuery = countQuery.Where("match_start_date = ?", params.Date)
	}
	if params.Status != "" {
		baseQuery = baseQuery.Where("match_status = ?", params.Status)
		countQuery = countQuery.Where("match_status = ?", params.Status)
	}
	if params.LeagueID != nil {
		baseQuery = baseQuery.Where("league_id = ?", *params.LeagueID)
		countQuery = countQuery.Where("league_id = ?", *params.LeagueID)
	}

	// Get total count
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var total int
	if err := db.Conn.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count matches: %w", err)
	}

	// Apply pagination and ordering
	baseQuery = baseQuery.
		OrderBy("match_start_date DESC", "match_start_time DESC").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset))

	sql, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var matches []SoccerMatch
	for rows.Next() {
		var m SoccerMatch
		err := rows.Scan(
			&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName,
			&m.MatchStatus, &m.MatchStartDate, &m.MatchStartTime,
			&m.HTeamID, &m.ATeamID, &m.HTeamName, &m.ATeamName,
			&m.HTeamGoals, &m.ATeamGoals, &m.HTScore, &m.FTScore,
			&m.Events, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		matches = append(matches, m)
	}

	return matches, total, nil
}

// GetLiveSoccerMatches returns currently live soccer matches
func (db *DB) GetLiveSoccerMatches() ([]SoccerMatch, error) {
	// Common statuses for live matches
	liveStatuses := []string{"1H", "HT", "2H", "ET", "P", "Live", "In Play"}

	query := db.Builder.
		Select("*").
		From("soccer_matches").
		Where("match_status IN (?, ?, ?, ?, ?, ?, ?)", liveStatuses[0], liveStatuses[1], liveStatuses[2], liveStatuses[3], liveStatuses[4], liveStatuses[5], liveStatuses[6]).
		OrderBy("match_start_time ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var matches []SoccerMatch
	for rows.Next() {
		var m SoccerMatch
		err := rows.Scan(
			&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName,
			&m.MatchStatus, &m.MatchStartDate, &m.MatchStartTime,
			&m.HTeamID, &m.ATeamID, &m.HTeamName, &m.ATeamName,
			&m.HTeamGoals, &m.ATeamGoals, &m.HTScore, &m.FTScore,
			&m.Events, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		matches = append(matches, m)
	}

	return matches, nil
}

// GetSoccerLeagues returns distinct leagues from soccer matches
func (db *DB) GetSoccerLeagues() ([]LeagueInfo, error) {
	query := db.Builder.
		Select("DISTINCT league_id", "league_gid", "league_name").
		From("soccer_matches").
		Where("league_id IS NOT NULL").
		OrderBy("league_name ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var leagues []LeagueInfo
	for rows.Next() {
		var l LeagueInfo
		var gid, id *int64
		var name *string
		err := rows.Scan(&id, &gid, &name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if id != nil {
			l.ID = *id
		}
		if gid != nil {
			l.GID = *gid
		}
		if name != nil {
			l.Name = *name
		}
		leagues = append(leagues, l)
	}

	return leagues, nil
}

// ============================================================================
// Basketball Filtered Queries (for API)
// ============================================================================

// GetBasketballMatchesFiltered returns basketball matches with filtering and pagination
func (db *DB) GetBasketballMatchesFiltered(params QueryParams) ([]BasketballMatch, int, error) {
	// Build base query
	baseQuery := db.Builder.
		Select("*").
		From("basketball_matches")

	countQuery := db.Builder.
		Select("COUNT(*)").
		From("basketball_matches")

	// Apply filters
	if params.Date != "" {
		baseQuery = baseQuery.Where("match_date = ?", params.Date)
		countQuery = countQuery.Where("match_date = ?", params.Date)
	}
	if params.Status != "" {
		baseQuery = baseQuery.Where("match_status = ?", params.Status)
		countQuery = countQuery.Where("match_status = ?", params.Status)
	}
	if params.LeagueID != nil {
		baseQuery = baseQuery.Where("league_id = ?", *params.LeagueID)
		countQuery = countQuery.Where("league_id = ?", *params.LeagueID)
	}

	// Get total count
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var total int
	if err := db.Conn.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count matches: %w", err)
	}

	// Apply pagination and ordering
	baseQuery = baseQuery.
		OrderBy("match_date DESC", "match_time DESC").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset))

	sql, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var matches []BasketballMatch
	for rows.Next() {
		var m BasketballMatch
		err := rows.Scan(
			&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName, &m.FileGroup,
			&m.MatchStatus, &m.MatchDate, &m.MatchTime, &m.Timer,
			&m.HTeamID, &m.HTeamName, &m.HTeamScore,
			&m.HTeamQ1, &m.HTeamQ2, &m.HTeamQ3, &m.HTeamQ4, &m.HTeamOt,
			&m.ATeamID, &m.ATeamName, &m.ATeamScore,
			&m.ATeamQ1, &m.ATeamQ2, &m.ATeamQ3, &m.ATeamQ4, &m.ATeamOt,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		matches = append(matches, m)
	}

	return matches, total, nil
}

// GetLiveBasketballMatches returns currently live basketball matches
func (db *DB) GetLiveBasketballMatches() ([]BasketballMatch, error) {
	// Common statuses for live basketball matches
	liveStatuses := []string{"Q1", "Q2", "Q3", "Q4", "OT", "HT", "Live", "In Play"}

	query := db.Builder.
		Select("*").
		From("basketball_matches").
		Where("match_status IN (?, ?, ?, ?, ?, ?, ?, ?)", liveStatuses[0], liveStatuses[1], liveStatuses[2], liveStatuses[3], liveStatuses[4], liveStatuses[5], liveStatuses[6], liveStatuses[7]).
		OrderBy("match_time ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var matches []BasketballMatch
	for rows.Next() {
		var m BasketballMatch
		err := rows.Scan(
			&m.ID, &m.MatchID, &m.LeagueGID, &m.LeagueID, &m.LeagueName, &m.FileGroup,
			&m.MatchStatus, &m.MatchDate, &m.MatchTime, &m.Timer,
			&m.HTeamID, &m.HTeamName, &m.HTeamScore,
			&m.HTeamQ1, &m.HTeamQ2, &m.HTeamQ3, &m.HTeamQ4, &m.HTeamOt,
			&m.ATeamID, &m.ATeamName, &m.ATeamScore,
			&m.ATeamQ1, &m.ATeamQ2, &m.ATeamQ3, &m.ATeamQ4, &m.ATeamOt,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		matches = append(matches, m)
	}

	return matches, nil
}

// GetBasketballLeagues returns distinct leagues from basketball matches
func (db *DB) GetBasketballLeagues() ([]LeagueInfo, error) {
	query := db.Builder.
		Select("DISTINCT league_id", "league_gid", "league_name").
		From("basketball_matches").
		Where("league_id IS NOT NULL").
		OrderBy("league_name ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var leagues []LeagueInfo
	for rows.Next() {
		var l LeagueInfo
		var gid, id *int64
		var name *string
		err := rows.Scan(&id, &gid, &name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if id != nil {
			l.ID = *id
		}
		if gid != nil {
			l.GID = *gid
		}
		if name != nil {
			l.Name = *name
		}
		leagues = append(leagues, l)
	}

	return leagues, nil
}

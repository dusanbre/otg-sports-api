package database

import (
	"fmt"
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

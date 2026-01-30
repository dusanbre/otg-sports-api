package database

import (
	"database/sql"
	"time"
)

// SoccerMatch represents a soccer match record
type SoccerMatch struct {
	ID             int64          `json:"id"`
	MatchID        sql.NullInt64  `json:"match_id"`
	LeagueGID      sql.NullInt64  `json:"league_gid"`
	LeagueID       sql.NullInt64  `json:"league_id"`
	LeagueName     sql.NullString `json:"league_name"`
	MatchStatus    sql.NullString `json:"match_status"`
	MatchStartDate sql.NullTime   `json:"match_start_date"`
	MatchStartTime sql.NullString `json:"match_start_time"`
	HTeamID        sql.NullInt64  `json:"h_team_id"`
	ATeamID        sql.NullInt64  `json:"a_team_id"`
	HTeamName      sql.NullString `json:"h_team_name"`
	ATeamName      sql.NullString `json:"a_team_name"`
	HTeamGoals     sql.NullInt32  `json:"h_team_goals"`
	ATeamGoals     sql.NullInt32  `json:"a_team_goals"`
	HTScore        sql.NullString `json:"ht_score"`
	FTScore        sql.NullString `json:"ft_score"`
	Events         sql.NullString `json:"events"` // JSON stored as string
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// BasketballMatch represents a basketball match record
type BasketballMatch struct {
	ID          int64          `json:"id"`
	MatchID     sql.NullInt64  `json:"match_id"`
	LeagueGID   sql.NullInt64  `json:"league_gid"`
	LeagueID    sql.NullInt64  `json:"league_id"`
	LeagueName  sql.NullString `json:"league_name"`
	FileGroup   sql.NullString `json:"file_group"`
	MatchStatus sql.NullString `json:"match_status"`
	MatchDate   sql.NullTime   `json:"match_date"`
	MatchTime   sql.NullString `json:"match_time"`
	Timer       sql.NullString `json:"timer"`
	HTeamID     sql.NullInt64  `json:"h_team_id"`
	HTeamName   sql.NullString `json:"h_team_name"`
	HTeamScore  sql.NullInt32  `json:"h_team_score"`
	HTeamQ1     sql.NullInt32  `json:"h_team_q1"`
	HTeamQ2     sql.NullInt32  `json:"h_team_q2"`
	HTeamQ3     sql.NullInt32  `json:"h_team_q3"`
	HTeamQ4     sql.NullInt32  `json:"h_team_q4"`
	HTeamOt     sql.NullInt32  `json:"h_team_ot"`
	ATeamID     sql.NullInt64  `json:"a_team_id"`
	ATeamName   sql.NullString `json:"a_team_name"`
	ATeamScore  sql.NullInt32  `json:"a_team_score"`
	ATeamQ1     sql.NullInt32  `json:"a_team_q1"`
	ATeamQ2     sql.NullInt32  `json:"a_team_q2"`
	ATeamQ3     sql.NullInt32  `json:"a_team_q3"`
	ATeamQ4     sql.NullInt32  `json:"a_team_q4"`
	ATeamOt     sql.NullInt32  `json:"a_team_ot"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ApiKey represents an API key record for authentication
type ApiKey struct {
	ID         int64        `json:"id"`
	KeyHash    string       `json:"key_hash"`
	KeyPrefix  string       `json:"key_prefix"`
	Name       string       `json:"name"`
	Sports     []string     `json:"sports"`     // Unmarshaled from JSON
	RateLimit  int          `json:"rate_limit"` // Requests per minute
	IsActive   bool         `json:"is_active"`
	CreatedAt  time.Time    `json:"created_at"`
	LastUsedAt sql.NullTime `json:"last_used_at"`
	ExpiresAt  sql.NullTime `json:"expires_at"`
}

// QueryParams holds common query parameters for filtering
type QueryParams struct {
	Limit    int
	Offset   int
	Date     string
	Status   string
	LeagueID *int64
}

// LeagueInfo represents league information
type LeagueInfo struct {
	ID   int64  `json:"id"`
	GID  int64  `json:"gid"`
	Name string `json:"name"`
}

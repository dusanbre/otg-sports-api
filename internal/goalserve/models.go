package goalserve

import (
	"encoding/json"
)

// GoalServeScores represents the root scores structure from GoalServe JSON API
type GoalServeScores struct {
	Categories []GoalServeScoreCategory `json:"category"`
}

// GoalServeScoreCategory represents a league/competition category
type GoalServeScoreCategory struct {
	ID      string               `json:"@id"`
	Gid     string               `json:"@gid"`
	Name    string               `json:"@name"`
	Matches GoalServeMatchesData `json:"matches"`
}

// GoalServeMatchesData wraps the match array/object
type GoalServeMatchesData struct {
	Match []GoalServeMatch
}

// UnmarshalJSON handles both single match object and array of matches
func (m *GoalServeMatchesData) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as an object first
	var temp struct {
		Match json.RawMessage `json:"match"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Check if match is an array or single object
	if len(temp.Match) > 0 && temp.Match[0] == '[' {
		// It's an array
		return json.Unmarshal(temp.Match, &m.Match)
	} else if len(temp.Match) > 0 && temp.Match[0] == '{' {
		// It's a single object, wrap it in an array
		var singleMatch GoalServeMatch
		if err := json.Unmarshal(temp.Match, &singleMatch); err != nil {
			return err
		}
		m.Match = []GoalServeMatch{singleMatch}
		return nil
	}

	// No matches
	m.Match = []GoalServeMatch{}
	return nil
}

// GoalServeMatch represents a soccer match from GoalServe
type GoalServeMatch struct {
	ID            string         `json:"@id"`
	Date          string         `json:"@date"`
	FormattedDate string         `json:"@formatted_date"`
	Time          string         `json:"@time"`
	Status        string         `json:"@status"`
	LocalTeam     GoalServeTeam  `json:"localteam"`
	VisitorTeam   GoalServeTeam  `json:"visitorteam"`
	HTScore       GoalServeScore `json:"ht"`
	FTScore       GoalServeScore `json:"ft"`
	Events        interface{}    `json:"events"` // Can be null or object with event array
}

// GoalServeTeam represents a team in a match
type GoalServeTeam struct {
	ID    string `json:"@id"`
	Name  string `json:"@name"`
	Goals string `json:"@goals"`
}

// GoalServeScore represents a score (halftime or fulltime)
type GoalServeScore struct {
	Score string `json:"@score"`
}

// GoalServeEvents wraps the event array
type GoalServeEvents struct {
	Event []GoalServeEvent `json:"event"`
}

// GoalServeEvent represents a match event
type GoalServeEvent struct {
	Type   string `json:"@type"`
	Team   string `json:"@team"`
	Player string `json:"@player"`
	Time   string `json:"@time"`
}

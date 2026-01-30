package goalserve

import (
	"encoding/json"
)

// GoalServeSoccerScores represents the root scores structure from GoalServe Soccer JSON API
type GoalServeSoccerScores struct {
	Categories []GoalServeSoccerCategory `json:"category"`
}

// GoalServeSoccerCategory represents a soccer league/competition category
type GoalServeSoccerCategory struct {
	ID      string                     `json:"@id"`
	Gid     string                     `json:"@gid"`
	Name    string                     `json:"@name"`
	Matches GoalServeSoccerMatchesData `json:"matches"`
}

// GoalServeSoccerMatchesData wraps the soccer match array/object
type GoalServeSoccerMatchesData struct {
	Match []GoalServeSoccerMatch
}

// UnmarshalJSON handles both single match object and array of matches
func (m *GoalServeSoccerMatchesData) UnmarshalJSON(data []byte) error {
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
		var singleMatch GoalServeSoccerMatch
		if err := json.Unmarshal(temp.Match, &singleMatch); err != nil {
			return err
		}
		m.Match = []GoalServeSoccerMatch{singleMatch}
		return nil
	}

	// No matches
	m.Match = []GoalServeSoccerMatch{}
	return nil
}

// GoalServeSoccerMatch represents a soccer match from GoalServe
type GoalServeSoccerMatch struct {
	ID            string               `json:"@id"`
	Date          string               `json:"@date"`
	FormattedDate string               `json:"@formatted_date"`
	Time          string               `json:"@time"`
	Status        string               `json:"@status"`
	LocalTeam     GoalServeSoccerTeam  `json:"localteam"`
	VisitorTeam   GoalServeSoccerTeam  `json:"visitorteam"`
	HTScore       GoalServeSoccerScore `json:"ht"`
	FTScore       GoalServeSoccerScore `json:"ft"`
	Events        interface{}          `json:"events"` // Can be null or object with event array
}

// GoalServeSoccerTeam represents a team in a soccer match
type GoalServeSoccerTeam struct {
	ID    string `json:"@id"`
	Name  string `json:"@name"`
	Goals string `json:"@goals"`
}

// GoalServeSoccerScore represents a soccer score (halftime or fulltime)
type GoalServeSoccerScore struct {
	Score string `json:"@score"`
}

// GoalServeSoccerEvents wraps the soccer event array
type GoalServeSoccerEvents struct {
	Event []GoalServeSoccerEvent `json:"event"`
}

// GoalServeSoccerEvent represents a soccer match event
type GoalServeSoccerEvent struct {
	Type   string `json:"@type"`
	Team   string `json:"@team"`
	Player string `json:"@player"`
	Time   string `json:"@time"`
}

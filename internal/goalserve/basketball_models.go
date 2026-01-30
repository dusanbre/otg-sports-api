package goalserve

import (
	"encoding/json"
)

// GoalServeBasketballScores represents the root scores structure from GoalServe Basketball JSON API
type GoalServeBasketballScores struct {
	Categories []GoalServeBasketballCategory `json:"category"`
}

// GoalServeBasketballCategory represents a basketball league/competition category
type GoalServeBasketballCategory struct {
	ID        string                       `json:"id"`
	Gid       string                       `json:"gid"`
	Name      string                       `json:"name"`
	FileGroup string                       `json:"file_group"`
	Match     GoalServeBasketballMatchData `json:"match"`
}

// GoalServeBasketballMatchData wraps the basketball match array/object
type GoalServeBasketballMatchData struct {
	Matches []GoalServeBasketballMatch
}

// UnmarshalJSON handles both single match object and array of matches for basketball
func (m *GoalServeBasketballMatchData) UnmarshalJSON(data []byte) error {
	// Check if it's an array or single object
	if len(data) > 0 && data[0] == '[' {
		// It's an array
		return json.Unmarshal(data, &m.Matches)
	} else if len(data) > 0 && data[0] == '{' {
		// It's a single object, wrap it in an array
		var singleMatch GoalServeBasketballMatch
		if err := json.Unmarshal(data, &singleMatch); err != nil {
			return err
		}
		m.Matches = []GoalServeBasketballMatch{singleMatch}
		return nil
	}

	// No matches or null
	m.Matches = []GoalServeBasketballMatch{}
	return nil
}

// GoalServeBasketballMatch represents a basketball match from GoalServe
type GoalServeBasketballMatch struct {
	ID        string                  `json:"id"`
	Date      string                  `json:"date"`
	Time      string                  `json:"time"`
	Status    string                  `json:"status"`
	Timer     string                  `json:"timer"`
	LocalTeam GoalServeBasketballTeam `json:"localteam"`
	AwayTeam  GoalServeBasketballTeam `json:"awayteam"`
}

// GoalServeBasketballTeam represents a team in a basketball match
type GoalServeBasketballTeam struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	TotalScore string `json:"totalscore"`
	Q1         string `json:"q1"`
	Q2         string `json:"q2"`
	Q3         string `json:"q3"`
	Q4         string `json:"q4"`
	Ot         string `json:"ot"`
}

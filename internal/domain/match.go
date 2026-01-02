package domain

import (
	"time"

	"github.com/google/uuid"
)

// Match representa un partido entre dos equipos
type Match struct {
	ID              uuid.UUID `json:"id"`
	MatchNumber     int       `json:"match_number"`
	Date            time.Time `json:"date"`
	Team1ID         uuid.UUID `json:"team1_id"`
	Team2ID         uuid.UUID `json:"team2_id"`
	GoalScoredTeam1 int       `json:"goal_scored_team1"`
	GoalScoredTeam2 int       `json:"goal_scored_team2"`
	CreatedAt       time.Time `json:"created_at"`
	// Relaciones opcionales
	Team1 *Team `json:"team1,omitempty"`
	Team2 *Team `json:"team2,omitempty"`
}

// NewMatch crea un nuevo partido
func NewMatch(matchNumber int, date time.Time, team1ID, team2ID uuid.UUID, goals1, goals2 int) *Match {
	return &Match{
		ID:              uuid.New(),
		MatchNumber:     matchNumber,
		Date:            date,
		Team1ID:         team1ID,
		Team2ID:         team2ID,
		GoalScoredTeam1: goals1,
		GoalScoredTeam2: goals2,
		CreatedAt:       time.Now().UTC(),
	}
}

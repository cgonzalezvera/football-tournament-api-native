package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tournament representa un torneo de fútbol
type Tournament struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	// Teams se carga bajo demanda
	Teams []Team `json:"teams,omitempty"`
}

// NewTournament crea un nuevo torneo
func NewTournament(name string) *Tournament {
	return &Tournament{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Teams:     []Team{},
	}
}

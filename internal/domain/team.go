package domain

import (
	"time"

	"github.com/google/uuid"
)

// Team representa un equipo de fútbol
type Team struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	// Players se carga bajo demanda, no siempre está presente
	Players []Player `json:"players,omitempty"`
}

// NewTeam crea un nuevo equipo
func NewTeam(name string) *Team {
	return &Team{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Players:   []Player{},
	}
}

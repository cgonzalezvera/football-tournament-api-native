package domain

import (
	"time"

	"github.com/google/uuid"
)

// Player representa un jugador de fútbol
// Equivalente a una entidad en C# con propiedades
type Player struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	DateBirth time.Time `json:"date_birth"`
	CreatedAt time.Time `json:"created_at"`
}

// NewPlayer crea un nuevo jugador con ID generado
// Equivalente a un constructor en C#
func NewPlayer(name string, dateBirth time.Time) *Player {
	return &Player{
		ID:        uuid.New(),
		Name:      name,
		DateBirth: dateBirth,
		CreatedAt: time.Now().UTC(),
	}
}

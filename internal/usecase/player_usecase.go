package usecase

import (
	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/repository"
	"github.com/google/uuid"
)

// PlayerUseCase contiene la lógica de negocio para jugadores
// Equivalente a un Service en C#
type PlayerUseCase struct {
	repo repository.PlayerRepository
}

func NewPlayerUseCase(repo repository.PlayerRepository) *PlayerUseCase {
	return &PlayerUseCase{repo: repo}
}

func (uc *PlayerUseCase) CreatePlayer(player *domain.Player) error {
	return uc.repo.Create(player)
}

func (uc *PlayerUseCase) GetPlayerByID(id uuid.UUID) (*domain.Player, error) {
	return uc.repo.GetByID(id)
}

func (uc *PlayerUseCase) GetAllPlayers() ([]domain.Player, error) {
	return uc.repo.GetAll()
}

func (uc *PlayerUseCase) UpdatePlayer(player *domain.Player) error {
	return uc.repo.Update(player)
}

func (uc *PlayerUseCase) DeletePlayer(id uuid.UUID) error {
	return uc.repo.Delete(id)
}

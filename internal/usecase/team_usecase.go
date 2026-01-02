package usecase

import (
	"fmt"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/repository"
	"github.com/google/uuid"
)

type TeamUseCase struct {
	teamRepo   repository.TeamRepository
	playerRepo repository.PlayerRepository
}

func NewTeamUseCase(teamRepo repository.TeamRepository, playerRepo repository.PlayerRepository) *TeamUseCase {
	return &TeamUseCase{
		teamRepo:   teamRepo,
		playerRepo: playerRepo,
	}
}

func (uc *TeamUseCase) CreateTeam(team *domain.Team) error {
	return uc.teamRepo.Create(team)
}

func (uc *TeamUseCase) GetTeamByID(id uuid.UUID) (*domain.Team, error) {
	return uc.teamRepo.GetByID(id)
}

func (uc *TeamUseCase) GetAllTeams() ([]domain.Team, error) {
	return uc.teamRepo.GetAll()
}

func (uc *TeamUseCase) UpdateTeam(team *domain.Team) error {
	return uc.teamRepo.Update(team)
}

func (uc *TeamUseCase) DeleteTeam(id uuid.UUID) error {
	return uc.teamRepo.Delete(id)
}

func (uc *TeamUseCase) AddPlayerToTeam(teamID, playerID uuid.UUID) error {
	// Validar que el equipo existe
	_, err := uc.teamRepo.GetByID(teamID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Validar que el jugador existe
	_, err = uc.playerRepo.GetByID(playerID)
	if err != nil {
		return fmt.Errorf("player not found: %w", err)
	}

	return uc.teamRepo.AddPlayer(teamID, playerID)
}

func (uc *TeamUseCase) RemovePlayerFromTeam(teamID, playerID uuid.UUID) error {
	return uc.teamRepo.RemovePlayer(teamID, playerID)
}

func (uc *TeamUseCase) GetTeamPlayers(teamID uuid.UUID) ([]domain.Player, error) {
	return uc.teamRepo.GetTeamPlayers(teamID)
}

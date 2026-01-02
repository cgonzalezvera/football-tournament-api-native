package usecase

import (
	"fmt"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/repository"
	"github.com/google/uuid"
)

type MatchUseCase struct {
	matchRepo repository.MatchRepository
	teamRepo  repository.TeamRepository
}

func NewMatchUseCase(matchRepo repository.MatchRepository, teamRepo repository.TeamRepository) *MatchUseCase {
	return &MatchUseCase{
		matchRepo: matchRepo,
		teamRepo:  teamRepo,
	}
}

func (uc *MatchUseCase) CreateMatch(match *domain.Match) error {
	// Validar que ambos equipos existen
	_, err := uc.teamRepo.GetByID(match.Team1ID)
	if err != nil {
		return fmt.Errorf("team1 not found: %w", err)
	}

	_, err = uc.teamRepo.GetByID(match.Team2ID)
	if err != nil {
		return fmt.Errorf("team2 not found: %w", err)
	}

	// Validar que no sea el mismo equipo
	if match.Team1ID == match.Team2ID {
		return fmt.Errorf("a team cannot play against itself")
	}

	return uc.matchRepo.Create(match)
}

func (uc *MatchUseCase) GetMatchByID(id uuid.UUID) (*domain.Match, error) {
	return uc.matchRepo.GetByID(id)
}

func (uc *MatchUseCase) GetAllMatches() ([]domain.Match, error) {
	return uc.matchRepo.GetAll()
}

func (uc *MatchUseCase) UpdateMatch(match *domain.Match) error {
	// Validar equipos
	_, err := uc.teamRepo.GetByID(match.Team1ID)
	if err != nil {
		return fmt.Errorf("team1 not found: %w", err)
	}

	_, err = uc.teamRepo.GetByID(match.Team2ID)
	if err != nil {
		return fmt.Errorf("team2 not found: %w", err)
	}

	if match.Team1ID == match.Team2ID {
		return fmt.Errorf("a team cannot play against itself")
	}

	return uc.matchRepo.Update(match)
}

func (uc *MatchUseCase) DeleteMatch(id uuid.UUID) error {
	return uc.matchRepo.Delete(id)
}

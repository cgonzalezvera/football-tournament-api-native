package usecase

import (
	"fmt"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/repository"
	"github.com/google/uuid"
)

type TournamentUseCase struct {
	tournamentRepo repository.TournamentRepository
	teamRepo       repository.TeamRepository
}

func NewTournamentUseCase(tournamentRepo repository.TournamentRepository, teamRepo repository.TeamRepository) *TournamentUseCase {
	return &TournamentUseCase{
		tournamentRepo: tournamentRepo,
		teamRepo:       teamRepo,
	}
}

func (uc *TournamentUseCase) CreateTournament(tournament *domain.Tournament) error {
	return uc.tournamentRepo.Create(tournament)
}

func (uc *TournamentUseCase) GetTournamentByID(id uuid.UUID) (*domain.Tournament, error) {
	return uc.tournamentRepo.GetByID(id)
}

func (uc *TournamentUseCase) GetAllTournaments() ([]domain.Tournament, error) {
	return uc.tournamentRepo.GetAll()
}

func (uc *TournamentUseCase) UpdateTournament(tournament *domain.Tournament) error {
	return uc.tournamentRepo.Update(tournament)
}

func (uc *TournamentUseCase) DeleteTournament(id uuid.UUID) error {
	return uc.tournamentRepo.Delete(id)
}

func (uc *TournamentUseCase) AddTeamToTournament(tournamentID, teamID uuid.UUID) error {
	// Validar que el torneo existe
	_, err := uc.tournamentRepo.GetByID(tournamentID)
	if err != nil {
		return fmt.Errorf("tournament not found: %w", err)
	}

	// Validar que el equipo existe
	_, err = uc.teamRepo.GetByID(teamID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	return uc.tournamentRepo.AddTeam(tournamentID, teamID)
}

func (uc *TournamentUseCase) RemoveTeamFromTournament(tournamentID, teamID uuid.UUID) error {
	return uc.tournamentRepo.RemoveTeam(tournamentID, teamID)
}

func (uc *TournamentUseCase) GetTournamentTeams(tournamentID uuid.UUID) ([]domain.Team, error) {
	return uc.tournamentRepo.GetTournamentTeams(tournamentID)
}

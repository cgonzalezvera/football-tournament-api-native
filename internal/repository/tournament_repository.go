package repository

import (
	"database/sql"
	"fmt"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/google/uuid"
)

type TournamentRepository interface {
	Create(tournament *domain.Tournament) error
	GetByID(id uuid.UUID) (*domain.Tournament, error)
	GetAll() ([]domain.Tournament, error)
	Update(tournament *domain.Tournament) error
	Delete(id uuid.UUID) error
	AddTeam(tournamentID, teamID uuid.UUID) error
	RemoveTeam(tournamentID, teamID uuid.UUID) error
	GetTournamentTeams(tournamentID uuid.UUID) ([]domain.Team, error)
}

type PostgresTournamentRepository struct {
	db *sql.DB
}

func NewPostgresTournamentRepository(db *sql.DB) TournamentRepository {
	return &PostgresTournamentRepository{db: db}
}

func (r *PostgresTournamentRepository) Create(tournament *domain.Tournament) error {
	query := `INSERT INTO tournaments (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, tournament.ID, tournament.Name, tournament.CreatedAt)
	return err
}

func (r *PostgresTournamentRepository) GetByID(id uuid.UUID) (*domain.Tournament, error) {
	query := `SELECT id, name, created_at FROM tournaments WHERE id = $1`
	var tournament domain.Tournament
	err := r.db.QueryRow(query, id).Scan(&tournament.ID, &tournament.Name, &tournament.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tournament not found")
	}
	if err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *PostgresTournamentRepository) GetAll() ([]domain.Tournament, error) {
	query := `SELECT id, name, created_at FROM tournaments ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []domain.Tournament
	for rows.Next() {
		var t domain.Tournament
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, rows.Err()
}

func (r *PostgresTournamentRepository) Update(tournament *domain.Tournament) error {
	query := `UPDATE tournaments SET name = $2 WHERE id = $1`
	result, err := r.db.Exec(query, tournament.ID, tournament.Name)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tournament not found")
	}
	return nil
}

func (r *PostgresTournamentRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM tournaments WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tournament not found")
	}
	return nil
}

func (r *PostgresTournamentRepository) AddTeam(tournamentID, teamID uuid.UUID) error {
	query := `INSERT INTO tournament_teams (tournament_id, team_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, tournamentID, teamID)
	return err
}

func (r *PostgresTournamentRepository) RemoveTeam(tournamentID, teamID uuid.UUID) error {
	query := `DELETE FROM tournament_teams WHERE tournament_id = $1 AND team_id = $2`
	_, err := r.db.Exec(query, tournamentID, teamID)
	return err
}

func (r *PostgresTournamentRepository) GetTournamentTeams(tournamentID uuid.UUID) ([]domain.Team, error) {
	query := `
		SELECT t.id, t.name, t.created_at
		FROM teams t
		INNER JOIN tournament_teams tt ON t.id = tt.team_id
		WHERE tt.tournament_id = $1
		ORDER BY t.name
	`
	rows, err := r.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, rows.Err()
}

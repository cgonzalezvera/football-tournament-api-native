package repository

import (
	"database/sql"
	"fmt"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/google/uuid"
)

type TeamRepository interface {
	Create(team *domain.Team) error
	GetByID(id uuid.UUID) (*domain.Team, error)
	GetAll() ([]domain.Team, error)
	Update(team *domain.Team) error
	Delete(id uuid.UUID) error
	AddPlayer(teamID, playerID uuid.UUID) error
	RemovePlayer(teamID, playerID uuid.UUID) error
	GetTeamPlayers(teamID uuid.UUID) ([]domain.Player, error)
}

type PostgresTeamRepository struct {
	db *sql.DB
}

func NewPostgresTeamRepository(db *sql.DB) TeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (r *PostgresTeamRepository) Create(team *domain.Team) error {
	query := `
		INSERT INTO teams (id, name, created_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, team.ID, team.Name, team.CreatedAt)
	return err
}

func (r *PostgresTeamRepository) GetByID(id uuid.UUID) (*domain.Team, error) {
	query := `
		SELECT id, name, created_at
		FROM teams
		WHERE id = $1
	`
	var team domain.Team
	err := r.db.QueryRow(query, id).Scan(&team.ID, &team.Name, &team.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *PostgresTeamRepository) GetAll() ([]domain.Team, error) {
	query := `SELECT id, name, created_at FROM teams ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
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

func (r *PostgresTeamRepository) Update(team *domain.Team) error {
	query := `UPDATE teams SET name = $2 WHERE id = $1`
	result, err := r.db.Exec(query, team.ID, team.Name)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}

func (r *PostgresTeamRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM teams WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}

func (r *PostgresTeamRepository) AddPlayer(teamID, playerID uuid.UUID) error {
	query := `INSERT INTO team_players (team_id, player_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, teamID, playerID)
	return err
}

func (r *PostgresTeamRepository) RemovePlayer(teamID, playerID uuid.UUID) error {
	query := `DELETE FROM team_players WHERE team_id = $1 AND player_id = $2`
	_, err := r.db.Exec(query, teamID, playerID)
	return err
}

func (r *PostgresTeamRepository) GetTeamPlayers(teamID uuid.UUID) ([]domain.Player, error) {
	query := `
		SELECT p.id, p.name, p.date_birth, p.created_at
		FROM players p
		INNER JOIN team_players tp ON p.id = tp.player_id
		WHERE tp.team_id = $1
		ORDER BY p.name
	`
	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []domain.Player
	for rows.Next() {
		var player domain.Player
		if err := rows.Scan(&player.ID, &player.Name, &player.DateBirth, &player.CreatedAt); err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, rows.Err()
}

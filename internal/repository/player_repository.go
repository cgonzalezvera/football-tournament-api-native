package repository

import (
	"database/sql"
	"fmt"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/google/uuid"
)

// PlayerRepository define el contrato para acceso a datos de jugadores
// En C# esto sería una interfaz IPlayerRepository
type PlayerRepository interface {
	Create(player *domain.Player) error
	GetByID(id uuid.UUID) (*domain.Player, error)
	GetAll() ([]domain.Player, error)
	Update(player *domain.Player) error
	Delete(id uuid.UUID) error
}

// PostgresPlayerRepository implementa PlayerRepository para PostgreSQL
type PostgresPlayerRepository struct {
	db *sql.DB
}

// NewPostgresPlayerRepository crea una nueva instancia del repositorio
func NewPostgresPlayerRepository(db *sql.DB) PlayerRepository {
	return &PostgresPlayerRepository{db: db}
}

func (r *PostgresPlayerRepository) Create(player *domain.Player) error {
	query := `
		INSERT INTO players (id, name, date_birth, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, player.ID, player.Name, player.DateBirth, player.CreatedAt)
	return err
}

func (r *PostgresPlayerRepository) GetByID(id uuid.UUID) (*domain.Player, error) {
	query := `
		SELECT id, name, date_birth, created_at
		FROM players
		WHERE id = $1
	`
	var player domain.Player
	err := r.db.QueryRow(query, id).Scan(
		&player.ID,
		&player.Name,
		&player.DateBirth,
		&player.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("player not found")
	}
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *PostgresPlayerRepository) GetAll() ([]domain.Player, error) {
	query := `
		SELECT id, name, date_birth, created_at
		FROM players
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
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

func (r *PostgresPlayerRepository) Update(player *domain.Player) error {
	query := `
		UPDATE players
		SET name = $2, date_birth = $3
		WHERE id = $1
	`
	result, err := r.db.Exec(query, player.ID, player.Name, player.DateBirth)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("player not found")
	}
	return nil
}

func (r *PostgresPlayerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM players WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("player not found")
	}
	return nil
}

package repository

import (
	"database/sql"
	"fmt"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/google/uuid"
)

type MatchRepository interface {
	Create(match *domain.Match) error
	GetByID(id uuid.UUID) (*domain.Match, error)
	GetAll() ([]domain.Match, error)
	Update(match *domain.Match) error
	Delete(id uuid.UUID) error
}

type PostgresMatchRepository struct {
	db *sql.DB
}

func NewPostgresMatchRepository(db *sql.DB) MatchRepository {
	return &PostgresMatchRepository{db: db}
}

func (r *PostgresMatchRepository) Create(match *domain.Match) error {
	query := `
		INSERT INTO matches (id, match_number, date, team1_id, team2_id, goal_scored_team1, goal_scored_team2, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		match.ID,
		match.MatchNumber,
		match.Date,
		match.Team1ID,
		match.Team2ID,
		match.GoalScoredTeam1,
		match.GoalScoredTeam2,
		match.CreatedAt,
	)
	return err
}

func (r *PostgresMatchRepository) GetByID(id uuid.UUID) (*domain.Match, error) {
	query := `
		SELECT id, match_number, date, team1_id, team2_id, goal_scored_team1, goal_scored_team2, created_at
		FROM matches
		WHERE id = $1
	`
	var match domain.Match
	err := r.db.QueryRow(query, id).Scan(
		&match.ID,
		&match.MatchNumber,
		&match.Date,
		&match.Team1ID,
		&match.Team2ID,
		&match.GoalScoredTeam1,
		&match.GoalScoredTeam2,
		&match.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("match not found")
	}
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *PostgresMatchRepository) GetAll() ([]domain.Match, error) {
	query := `
		SELECT id, match_number, date, team1_id, team2_id, goal_scored_team1, goal_scored_team2, created_at
		FROM matches
		ORDER BY date DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		var match domain.Match
		if err := rows.Scan(
			&match.ID,
			&match.MatchNumber,
			&match.Date,
			&match.Team1ID,
			&match.Team2ID,
			&match.GoalScoredTeam1,
			&match.GoalScoredTeam2,
			&match.CreatedAt,
		); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, rows.Err()
}

func (r *PostgresMatchRepository) Update(match *domain.Match) error {
	query := `
		UPDATE matches
		SET match_number = $2, date = $3, team1_id = $4, team2_id = $5, 
		    goal_scored_team1 = $6, goal_scored_team2 = $7
		WHERE id = $1
	`
	result, err := r.db.Exec(query,
		match.ID,
		match.MatchNumber,
		match.Date,
		match.Team1ID,
		match.Team2ID,
		match.GoalScoredTeam1,
		match.GoalScoredTeam2,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("match not found")
	}
	return nil
}

func (r *PostgresMatchRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM matches WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("match not found")
	}
	return nil
}

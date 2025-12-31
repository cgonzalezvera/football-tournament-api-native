-- Schema inicial para la API de Torneo de Fútbol

-- Habilitar extensión para UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabla de Jugadores
CREATE TABLE IF NOT EXISTS players (
                                       id UUID PRIMARY KEY,
                                       name VARCHAR(255) NOT NULL,
    date_birth TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

-- Tabla de Equipos
CREATE TABLE IF NOT EXISTS teams (
                                     id UUID PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

-- Tabla de Torneos
CREATE TABLE IF NOT EXISTS tournaments (
                                           id UUID PRIMARY KEY,
                                           name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

-- Tabla de Partidos
CREATE TABLE IF NOT EXISTS matches (
                                       id UUID PRIMARY KEY,
                                       match_number INTEGER NOT NULL,
                                       date TIMESTAMP WITH TIME ZONE NOT NULL,
                                       team1_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    team2_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    goal_scored_team1 INTEGER NOT NULL DEFAULT 0,
    goal_scored_team2 INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT different_teams CHECK (team1_id != team2_id)
    );

-- Tabla de relación Team-Player (muchos a muchos)
CREATE TABLE IF NOT EXISTS team_players (
                                            team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (team_id, player_id)
    );

-- Tabla de relación Tournament-Team (muchos a muchos)
CREATE TABLE IF NOT EXISTS tournament_teams (
                                                tournament_id UUID NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tournament_id, team_id)
    );

-- Índices para mejorar el rendimiento de las consultas
CREATE INDEX IF NOT EXISTS idx_players_name ON players(name);
CREATE INDEX IF NOT EXISTS idx_players_date_birth ON players(date_birth);
CREATE INDEX IF NOT EXISTS idx_teams_name ON teams(name);
CREATE INDEX IF NOT EXISTS idx_matches_date ON matches(date);
CREATE INDEX IF NOT EXISTS idx_matches_team1 ON matches(team1_id);
CREATE INDEX IF NOT EXISTS idx_matches_team2 ON matches(team2_id);
CREATE INDEX IF NOT EXISTS idx_team_players_team ON team_players(team_id);
CREATE INDEX IF NOT EXISTS idx_team_players_player ON team_players(player_id);
CREATE INDEX IF NOT EXISTS idx_tournament_teams_tournament ON tournament_teams(tournament_id);
CREATE INDEX IF NOT EXISTS idx_tournament_teams_team ON tournament_teams(team_id);

-- Datos de ejemplo (opcional, para testing)
-- Puedes descomentar esto para tener datos iniciales

-- INSERT INTO players (id, name, date_birth, created_at) VALUES
-- (uuid_generate_v4(), 'Lionel Messi', '1987-06-24 00:00:00+00', NOW()),
-- (uuid_generate_v4(), 'Cristiano Ronaldo', '1985-02-05 00:00:00+00', NOW()),
-- (uuid_generate_v4(), 'Neymar Jr', '1992-02-05 00:00:00+00', NOW());

-- INSERT INTO teams (id, name, created_at) VALUES
-- (uuid_generate_v4(), 'Argentina', NOW()),
-- (uuid_generate_v4(), 'Portugal', NOW()),
-- (uuid_generate_v4(), 'Brasil', NOW());

-- Comentarios en las tablas para documentación
COMMENT ON TABLE players IS 'Almacena la información de los jugadores';
COMMENT ON TABLE teams IS 'Almacena la información de los equipos';
COMMENT ON TABLE tournaments IS 'Almacena la información de los torneos';
COMMENT ON TABLE matches IS 'Almacena la información de los partidos entre equipos';
COMMENT ON TABLE team_players IS 'Tabla de relación muchos a muchos entre equipos y jugadores';
COMMENT ON TABLE tournament_teams IS 'Tabla de relación muchos a muchos entre torneos y equipos';

-- Verificar que todo se creó correctamente
DO $$
BEGIN
    RAISE NOTICE 'Database schema created successfully!';
    RAISE NOTICE 'Tables: players, teams, tournaments, matches, team_players, tournament_teams';
END $$;
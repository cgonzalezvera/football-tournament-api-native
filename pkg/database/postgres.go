// pkg/database/postgres.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Driver de PostgreSQL
)

// Config contiene la configuración de conexión a la base de datos
// En C# esto sería similar a ConnectionStrings en appsettings.json
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConfigFromEnv crea una configuración desde variables de entorno
func NewConfigFromEnv() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "tournament_db"),
	}
}

// NewConnection crea una nueva conexión a PostgreSQL
// Equivalente a DbContext en Entity Framework
func NewConnection(config *Config) (*sql.DB, error) {
	// String de conexión de PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	// Abrir conexión
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)                 // Máximo de conexiones abiertas
	db.SetMaxIdleConns(5)                  // Conexiones en idle
	db.SetConnMaxLifetime(5 * time.Minute) // Tiempo de vida de conexión

	// Verificar conexión con timeout
	if err := pingWithRetry(db, 5); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL database")
	return db, nil
}

// pingWithRetry intenta conectarse a la base de datos con reintentos
func pingWithRetry(db *sql.DB, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}

		log.Printf("⏳ Waiting for database connection... (attempt %d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

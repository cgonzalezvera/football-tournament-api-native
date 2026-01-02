package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/handler"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/repository"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/usecase"
	"github.com/cgonzalezvera/football-tournament-api-native/pkg/database"
)

func main() {
	// Configurar logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("🚀 Starting Tournament API...")

	// Conectar a la base de datos
	dbConfig := database.NewConfigFromEnv()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inicializar repositorios (Data Access Layer)
	playerRepo := repository.NewPostgresPlayerRepository(db)
	teamRepo := repository.NewPostgresTeamRepository(db)
	tournamentRepo := repository.NewPostgresTournamentRepository(db)
	matchRepo := repository.NewPostgresMatchRepository(db)

	// Inicializar casos de uso (Business Logic Layer)
	playerUC := usecase.NewPlayerUseCase(playerRepo)
	teamUC := usecase.NewTeamUseCase(teamRepo, playerRepo)
	tournamentUC := usecase.NewTournamentUseCase(tournamentRepo, teamRepo)
	matchUC := usecase.NewMatchUseCase(matchRepo, teamRepo)

	// Inicializar handlers (Presentation Layer)
	playerHandler := handler.NewPlayerHandler(playerUC)
	teamHandler := handler.NewTeamHandler(teamUC)
	tournamentHandler := handler.NewTournamentHandler(tournamentUC)
	matchHandler := handler.NewMatchHandler(matchUC)

	// Configurar rutas (equivalente a app.MapControllers() en C#)
	mux := http.NewServeMux()

	// Rutas de jugadores
	mux.Handle("/api/players", enableCORS(playerHandler))
	mux.Handle("/api/players/", enableCORS(playerHandler))

	// Rutas de equipos
	mux.Handle("/api/teams", enableCORS(teamHandler))
	mux.Handle("/api/teams/", enableCORS(teamHandler))

	// Rutas de torneos
	mux.Handle("/api/tournaments", enableCORS(tournamentHandler))
	mux.Handle("/api/tournaments/", enableCORS(tournamentHandler))

	// Rutas de partidos
	mux.Handle("/api/matches", enableCORS(matchHandler))
	mux.Handle("/api/matches/", enableCORS(matchHandler))

	// Ruta de health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"tournament-api"}`))
	})

	// Obtener puerto desde variable de entorno
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor HTTP
	serverAddr := ":" + port
	log.Printf("🌐 Server listening on http://localhost%s", serverAddr)
	log.Printf("📚 Health check: http://localhost%s/health", serverAddr)
	log.Printf("📋 API Base URL: http://localhost%s/api", serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// enableCORS es un middleware para habilitar CORS
// En C# esto sería similar a app.UseCors() en Program.cs
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

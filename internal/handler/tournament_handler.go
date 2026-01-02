package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/usecase"
	"github.com/google/uuid"
)

type TournamentHandler struct {
	useCase *usecase.TournamentUseCase
}

func NewTournamentHandler(useCase *usecase.TournamentUseCase) *TournamentHandler {
	return &TournamentHandler{useCase: useCase}
}

func (h *TournamentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/tournaments")
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	// Manejar /api/tournaments/{id}/teams/{teamId}
	if len(segments) >= 3 && segments[1] == "teams" {
		tournamentID, err := uuid.Parse(segments[0])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid tournament UUID")
			return
		}

		teamID, err := uuid.Parse(segments[2])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid team UUID")
			return
		}

		switch r.Method {
		case http.MethodPost:
			h.AddTeam(w, r, tournamentID, teamID)
		case http.MethodDelete:
			h.RemoveTeam(w, r, tournamentID, teamID)
		default:
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// Manejar /api/tournaments/{id}/teams
	if len(segments) == 2 && segments[1] == "teams" {
		tournamentID, err := uuid.Parse(segments[0])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid tournament UUID")
			return
		}

		if r.Method == http.MethodGet {
			h.GetTournamentTeams(w, r, tournamentID)
			return
		}
	}

	// Rutas CRUD estándar
	switch r.Method {
	case http.MethodGet:
		if path == "" {
			h.GetAll(w, r)
		} else {
			h.GetByID(w, r, path)
		}
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodPut:
		h.Update(w, r, path)
	case http.MethodDelete:
		h.Delete(w, r, path)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TournamentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tournament := domain.NewTournament(input.Name)
	if err := h.useCase.CreateTournament(tournament); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, tournament)
}

func (h *TournamentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tournaments, err := h.useCase.GetAllTournaments()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tournaments)
}

func (h *TournamentHandler) GetByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	tournament, err := h.useCase.GetTournamentByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tournament)
}

func (h *TournamentHandler) Update(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tournament := &domain.Tournament{ID: id, Name: input.Name}
	if err := h.useCase.UpdateTournament(tournament); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tournament)
}

func (h *TournamentHandler) Delete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.useCase.DeleteTournament(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Tournament deleted"})
}

func (h *TournamentHandler) AddTeam(w http.ResponseWriter, r *http.Request, tournamentID, teamID uuid.UUID) {
	if err := h.useCase.AddTeamToTournament(tournamentID, teamID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Team added to tournament"})
}

func (h *TournamentHandler) RemoveTeam(w http.ResponseWriter, r *http.Request, tournamentID, teamID uuid.UUID) {
	if err := h.useCase.RemoveTeamFromTournament(tournamentID, teamID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Team removed from tournament"})
}

func (h *TournamentHandler) GetTournamentTeams(w http.ResponseWriter, r *http.Request, tournamentID uuid.UUID) {
	teams, err := h.useCase.GetTournamentTeams(tournamentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, teams)
}

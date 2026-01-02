package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/usecase"
	"github.com/google/uuid"
)

type TeamHandler struct {
	useCase *usecase.TeamUseCase
}

func NewTeamHandler(useCase *usecase.TeamUseCase) *TeamHandler {
	return &TeamHandler{useCase: useCase}
}

func (h *TeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/teams")
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	// Manejar rutas como /api/teams/{id}/players/{playerId}
	if len(segments) >= 3 && segments[1] == "players" {
		teamID, err := uuid.Parse(segments[0])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid team UUID")
			return
		}

		playerID, err := uuid.Parse(segments[2])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid player UUID")
			return
		}

		switch r.Method {
		case http.MethodPost:
			h.AddPlayer(w, r, teamID, playerID)
		case http.MethodDelete:
			h.RemovePlayer(w, r, teamID, playerID)
		default:
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// Manejar rutas como /api/teams/{id}/players
	if len(segments) == 2 && segments[1] == "players" {
		teamID, err := uuid.Parse(segments[0])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid team UUID")
			return
		}

		if r.Method == http.MethodGet {
			h.GetTeamPlayers(w, r, teamID)
			return
		}
	}

	// Rutas estándar CRUD
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

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	team := domain.NewTeam(input.Name)
	if err := h.useCase.CreateTeam(team); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, team)
}

func (h *TeamHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	teams, err := h.useCase.GetAllTeams()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, teams)
}

func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	team, err := h.useCase.GetTeamByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request, idStr string) {
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

	team := &domain.Team{ID: id, Name: input.Name}
	if err := h.useCase.UpdateTeam(team); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.useCase.DeleteTeam(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Team deleted"})
}

func (h *TeamHandler) AddPlayer(w http.ResponseWriter, r *http.Request, teamID, playerID uuid.UUID) {
	if err := h.useCase.AddPlayerToTeam(teamID, playerID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Player added to team"})
}

func (h *TeamHandler) RemovePlayer(w http.ResponseWriter, r *http.Request, teamID, playerID uuid.UUID) {
	if err := h.useCase.RemovePlayerFromTeam(teamID, playerID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Player removed from team"})
}

func (h *TeamHandler) GetTeamPlayers(w http.ResponseWriter, r *http.Request, teamID uuid.UUID) {
	players, err := h.useCase.GetTeamPlayers(teamID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, players)
}

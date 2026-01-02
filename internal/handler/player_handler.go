package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/usecase"
	"github.com/google/uuid"
)

type PlayerHandler struct {
	useCase *usecase.PlayerUseCase
}

func NewPlayerHandler(useCase *usecase.PlayerUseCase) *PlayerHandler {
	return &PlayerHandler{useCase: useCase}
}

// En Go no hay atributos como [HttpGet], usamos funciones que verifican el método
func (h *PlayerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la URL si existe: /api/players/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/players")
	path = strings.Trim(path, "/")

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

func (h *PlayerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		DateBirth string `json:"date_birth"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dateBirth, err := parseDateTime(input.DateBirth)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format, use ISO 8601")
		return
	}

	player := domain.NewPlayer(input.Name, dateBirth)
	if err := h.useCase.CreatePlayer(player); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, player)
}

func (h *PlayerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	players, err := h.useCase.GetAllPlayers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, players)
}

func (h *PlayerHandler) GetByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	player, err := h.useCase.GetPlayerByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}

func (h *PlayerHandler) Update(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var input struct {
		Name      string `json:"name"`
		DateBirth string `json:"date_birth"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dateBirth, err := parseDateTime(input.DateBirth)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	player := &domain.Player{
		ID:        id,
		Name:      input.Name,
		DateBirth: dateBirth,
	}

	if err := h.useCase.UpdatePlayer(player); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}

func (h *PlayerHandler) Delete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.useCase.DeletePlayer(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Player deleted"})
}

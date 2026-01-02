package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cgonzalezvera/football-tournament-api-native/internal/domain"
	"github.com/cgonzalezvera/football-tournament-api-native/internal/usecase"
	"github.com/google/uuid"
)

type MatchHandler struct {
	useCase *usecase.MatchUseCase
}

func NewMatchHandler(useCase *usecase.MatchUseCase) *MatchHandler {
	return &MatchHandler{useCase: useCase}
}

func (h *MatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/matches")
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

func (h *MatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		MatchNumber     int    `json:"match_number"`
		Date            string `json:"date"`
		Team1ID         string `json:"team1_id"`
		Team2ID         string `json:"team2_id"`
		GoalScoredTeam1 int    `json:"goal_scored_team1"`
		GoalScoredTeam2 int    `json:"goal_scored_team2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	date, err := parseDateTime(input.Date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	team1ID, err := uuid.Parse(input.Team1ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team1_id UUID")
		return
	}

	team2ID, err := uuid.Parse(input.Team2ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team2_id UUID")
		return
	}

	match := domain.NewMatch(
		input.MatchNumber,
		date,
		team1ID,
		team2ID,
		input.GoalScoredTeam1,
		input.GoalScoredTeam2,
	)

	if err := h.useCase.CreateMatch(match); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, match)
}

func (h *MatchHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	matches, err := h.useCase.GetAllMatches()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, matches)
}

func (h *MatchHandler) GetByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	match, err := h.useCase.GetMatchByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, match)
}

func (h *MatchHandler) Update(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var input struct {
		MatchNumber     int    `json:"match_number"`
		Date            string `json:"date"`
		Team1ID         string `json:"team1_id"`
		Team2ID         string `json:"team2_id"`
		GoalScoredTeam1 int    `json:"goal_scored_team1"`
		GoalScoredTeam2 int    `json:"goal_scored_team2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	date, err := parseDateTime(input.Date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	team1ID, err := uuid.Parse(input.Team1ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team1_id UUID")
		return
	}

	team2ID, err := uuid.Parse(input.Team2ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team2_id UUID")
		return
	}

	match := &domain.Match{
		ID:              id,
		MatchNumber:     input.MatchNumber,
		Date:            date,
		Team1ID:         team1ID,
		Team2ID:         team2ID,
		GoalScoredTeam1: input.GoalScoredTeam1,
		GoalScoredTeam2: input.GoalScoredTeam2,
	}

	if err := h.useCase.UpdateMatch(match); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, match)
}

func (h *MatchHandler) Delete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.useCase.DeleteMatch(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Match deleted"})
}

package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Funciones helper para respuestas HTTP (equivalente a ActionResult en C#)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func parseDateTime(dateStr string) (time.Time, error) {
	// Parsear fecha en formato RFC3339 (ISO 8601)
	// Ejemplo: "2023-06-24T00:00:00Z"
	return time.Parse(time.RFC3339, dateStr)
}

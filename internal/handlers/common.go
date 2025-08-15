package handlers

import (
	"encoding/json"
	"net/http"

	"golang/tutorial/todo/internal/validation"
)

// NOTE: ヘルパー関数のためanyを許容
func WriteJSON(w http.ResponseWriter, statusCode int, v ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(v) > 0 {
		_ = json.NewEncoder(w).Encode(v[0])
	}
}

func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func ReadJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}

func WriteValidationError(w http.ResponseWriter, vErr validation.Errors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(vErr)
}

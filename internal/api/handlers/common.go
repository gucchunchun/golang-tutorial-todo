package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang/tutorial/todo/internal/apperr"
)

// NOTE: ヘルパー関数のためanyを許容
func WriteJSON(w http.ResponseWriter, statusCode int, v ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(v) > 0 {
		_ = json.NewEncoder(w).Encode(v[0])
	}
}

type ErrorJSON struct {
	Code    string
	Message string
}

func writeErrorJSON(w http.ResponseWriter, statusCode int, errCode string, message string) {
	WriteJSON(w, statusCode, ErrorJSON{Code: errCode, Message: message})
}
func WriteError(w http.ResponseWriter, err error) {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		writeErrorJSON(w, ae.Code.HTTPStatus(), ae.Code.String(), ae.Error())
		return
	}
	writeErrorJSON(w, http.StatusInternalServerError, apperr.CodeUnknown.String(), err.Error())
}

func ReadJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}

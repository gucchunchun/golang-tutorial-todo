package httpapi

import (
	"encoding/json"
	"net/http"

	"golang/tutorial/todo/internal/models"
)

type TaskService interface {
	// Create(title string) error
	// Get(id string) (models.Task, error)
	List() ([]models.Task, error)
	// Update(id string, title string) error
	// Delete(id string) error
}

type Server struct {
	// TaskService TaskService
}

func New() *Server {
	return &Server{
		// TaskService: taskService,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleHelloworld)
	return mux
}

// NOTE: ヘルパー関数のためanyを許容
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (s *Server) handleHelloworld(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Hello, World!"})
}

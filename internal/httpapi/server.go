package httpapi

import (
	"encoding/json"
	"net/http"
)

// type TaskHandler intercace

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

func (s *Server) handleHelloworld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

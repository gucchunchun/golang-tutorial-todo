package httpapi

import (
	"encoding/json"
	"net/http"

	"golang/tutorial/todo/internal/handlers/task"
)

type Handler interface {
	Routes(mux *http.ServeMux)
}

type Server struct {
	handlers []Handler
}

func New(taskService task.TaskService) *Server {
	return &Server{
		handlers: []Handler{task.NewTaskHandler(taskService)},
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleHelloworld)
	// ルーティングの設定
	for _, h := range s.handlers {
		h.Routes(mux)
	}
	return mux
}

func (s *Server) handleHelloworld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

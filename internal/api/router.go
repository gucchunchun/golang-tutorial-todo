package api

import (
	"encoding/json"
	"net/http"

	"golang/tutorial/todo/internal/handlers/task"
)

type Handler interface {
	Routes(mux *http.ServeMux)
}

type Router struct {
	handlers []Handler
}

func New(taskService task.TaskService) *Router {
	return &Router{
		handlers: []Handler{task.NewTaskHandler(taskService)},
	}
}

func (s *Router) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleHelloworld)
	// ルーティングの設定
	for _, h := range s.handlers {
		h.Routes(mux)
	}
	return mux
}

func (s *Router) handleHelloworld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

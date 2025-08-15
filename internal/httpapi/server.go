package httpapi

import (
	"context"
	"net/http"

	"golang/tutorial/todo/internal/models"
)

type Task interface {
	Create(ctx context.Context, title string) error
	Get(ctx context.Context, id string) (models.Task, error)
	List(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, id string, title string) error
	Delete(ctx context.Context, id string) error
}

type Server struct {
	// task Task
}

func New() *Server {
	return &Server{
		// task: task,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleHelloworld)
	return mux
}

func (s *Server) handleHelloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

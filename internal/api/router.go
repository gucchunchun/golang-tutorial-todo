package api

import (
	"encoding/json"
	"net/http"

	"golang/tutorial/todo/internal/api/handlers/task"
	"golang/tutorial/todo/internal/api/middleware/logmw"
	"golang/tutorial/todo/internal/api/middleware/recovermw"
)

type Middleware = func(next http.Handler) http.Handler

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

	// 個別にミドルウェア設定（この場合二度ログ出力される）
	mux.Handle("GET /", logmw.Log(http.HandlerFunc((s.handleHelloworld))))
	mux.HandleFunc("GET /health", (s.handleHealth))
	mux.HandleFunc("GET /panic", (s.handlePanic))

	// ルーティングの設定
	for _, h := range s.handlers {
		h.Routes(mux)
	}
	// 全てのルートにミドルウェアの設定
	middlewares := []Middleware{
		recovermw.Recover,
		logmw.Log,
	}
	var handler http.Handler = mux
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func (s *Router) handleHelloworld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

func (s *Router) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Router) handlePanic(w http.ResponseWriter, r *http.Request) {
	panic("panic")
}

package quotehdl

import (
	"context"
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
)

type QuoteService interface {
	GetRandomQuote(ctx context.Context) (string, error)
}

type QuoteHandler struct {
	QuoteService QuoteService
}

func New(quoteService QuoteService) *QuoteHandler {
	return &QuoteHandler{
		QuoteService: quoteService,
	}
}

func (h *QuoteHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /quote", h.get)
}

func (h *QuoteHandler) get(w http.ResponseWriter, r *http.Request) {
	q, err := h.QuoteService.GetRandomQuote(r.Context())
	if err != nil {
		handlers.WriteError(w, err)
		return
	}
	handlers.WriteJSON(w, http.StatusOK, q)
}

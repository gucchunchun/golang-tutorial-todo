package quotehdl

import (
	"context"
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
	"golang/tutorial/todo/internal/quote"
)

type QuoteService interface {
	RandomQuote(ctx context.Context) (quote.Quote, error)
}

type QuoteHandler struct {
	QuoteService QuoteService
}

func NewQuoteHandler(quoteService QuoteService) *QuoteHandler {
	return &QuoteHandler{
		QuoteService: quoteService,
	}
}

func (h *QuoteHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /quote", h.get)
}

func (h *QuoteHandler) get(w http.ResponseWriter, r *http.Request) {
	q, err := h.QuoteService.RandomQuote(r.Context())
	if err != nil {
		handlers.WriteError(w, err)
		return
	}
	handlers.WriteJSON(w, http.StatusOK, q)
}

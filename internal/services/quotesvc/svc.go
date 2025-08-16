package quotesvc

import (
	"context"
	"golang/tutorial/todo/internal/quotes"
)

type QuoteClient interface {
	RandomQuote(ctx context.Context) (quotes.Quote, error)
}

type QuoteService struct {
	quoteClient QuoteClient
}

func NewQuoteService(quoteClient QuoteClient) *QuoteService {
	return &QuoteService{
		quoteClient: quoteClient,
	}
}

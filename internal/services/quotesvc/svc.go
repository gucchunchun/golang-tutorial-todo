package quotesvc

import (
	"context"
	"golang/tutorial/todo/internal/quote"
)

type QuoteClient interface {
	RandomQuote(ctx context.Context) (quote.Quote, error)
}

type QuoteService struct {
	quoteClient QuoteClient
}

func NewQuoteService(quoteClient QuoteClient) *QuoteService {
	return &QuoteService{
		quoteClient: quoteClient,
	}
}

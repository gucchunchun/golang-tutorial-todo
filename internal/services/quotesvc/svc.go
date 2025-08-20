package quotesvc

import (
	"context"
	"golang/tutorial/todo/internal/quote"
)

type QuoteClient interface {
	RandomQuote(ctx context.Context) (quote.Quote, error)
}

type Service struct {
	quoteClient QuoteClient
}

func New(quoteClient QuoteClient) *Service {
	return &Service{
		quoteClient: quoteClient,
	}
}

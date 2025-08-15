package task

import (
	"context"
	"golang/tutorial/todo/internal/quotes"
	"golang/tutorial/todo/internal/storage"
)

type QuoteClient interface {
	RandomQuote(ctx context.Context) (quotes.Quote, error)
}

type Service struct {
	quoteClient QuoteClient
	storage     storage.Client
}

func NewService(quoteClient QuoteClient, storage storage.Client) *Service {
	return &Service{
		quoteClient: quoteClient,
		storage:     storage,
	}
}

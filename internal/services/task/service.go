package task

import (
	"context"
	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/quotes"
)

type QuoteClient interface {
	RandomQuote(ctx context.Context) (quotes.Quote, error)
}
type StorageClient interface {
	LoadTasks() ([]models.Task, error)
	SaveTasks(tasks []models.Task) error
}

type Service struct {
	quoteClient QuoteClient
	storage     StorageClient
}

func NewService(quoteClient QuoteClient, storage StorageClient) *Service {
	return &Service{
		quoteClient: quoteClient,
		storage:     storage,
	}
}

package task

import (
	"golang/tutorial/todo/internal/quotes"
	"golang/tutorial/todo/internal/storage"
)

type Service struct {
	quoteClient quotes.Client
	storage     storage.Client
}

func NewService(quoteClient quotes.Client, storage storage.Client) *Service {
	return &Service{
		quoteClient: quoteClient,
		storage:     storage,
	}
}

package quotesvc

import (
	"context"
	"fmt"

	"golang/tutorial/todo/internal/apperr"
)

func (s Service) GetRandomQuote(ctx context.Context) (string, error) {
	quote, err := s.quoteClient.RandomQuote(ctx)
	if err != nil {
		return "", apperr.E(apperr.CodeUnknown, "Failed to get random quote", err)
	}
	return fmt.Sprintf("%s - %s", quote.Author, quote.Text), nil
}

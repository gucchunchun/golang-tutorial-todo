package quotesvc

import (
	"context"
	"fmt"

	"golang/tutorial/todo/internal/apperr"
)

func (s *QuoteService) GetRandomQuote() (string, error) {
	quote, err := s.quoteClient.RandomQuote(context.Background())
	if err != nil {
		return "", apperr.E(apperr.CodeUnknown, "Failed to get random quote", err)
	}
	return fmt.Sprintf("%s - %s", quote.Author, quote.Text), nil
}

package task

import (
	"context"
	"golang/tutorial/todo/internal/quotes"
)

func GetQuote(c quotes.Client) (string, error) {
	quote, err := c.RandomQuote(context.Background())
	if err != nil {
		return "", err
	}
	return quote.Text, nil
}

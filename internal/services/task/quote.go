package task

import (
	"context"
)

func (c *Service) GetQuote(ctx context.Context) (string, error) {
	quote, err := c.quoteClient.RandomQuote(ctx)
	if err != nil {
		return "", err
	}
	return quote.Text, nil
}

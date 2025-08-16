package quotesvc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/quotes"
	"golang/tutorial/todo/internal/services/quotesvc"
)

// stub
type quoteClientFunc func(ctx context.Context) (quotes.Quote, error)

func (f quoteClientFunc) RandomQuote(ctx context.Context) (quotes.Quote, error) {
	return f(ctx)
}

var errStub = errors.New("stub")

func TestQuoteService(t *testing.T) {
	t.Run("GetRandomQuote", func(t *testing.T) {
		tests := map[string]struct {
			stubFunc func(ctx context.Context) (quotes.Quote, error)
			wantErr  bool
		}{
			"ok": {
				stubFunc: func(ctx context.Context) (quotes.Quote, error) {
					return quotes.Quote{Author: "Ada", Text: "Keep going"}, nil
				},
				wantErr: false,
			},
			"error": {
				stubFunc: func(ctx context.Context) (quotes.Quote, error) {
					return quotes.Quote{}, errStub
				},
				wantErr: true,
			},
		}

		for name, tc := range tests {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				svc := quotesvc.NewQuoteService(quoteClientFunc(tc.stubFunc))

				_, err := svc.GetRandomQuote()
				if !tc.wantErr {
					assert.NoError(t, err)
					return
				}

				assert.Error(t, err)
				var ae *apperr.Error
				assert.ErrorAs(t, err, &ae)
				assert.Equal(t, ae.Code, apperr.CodeUnknown)
				assert.Equal(t, ae.Message, "Failed to get random quote")
				assert.ErrorIs(t, ae.Err, errStub)
			})
		}
	})

}

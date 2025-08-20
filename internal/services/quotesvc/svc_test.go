package quotesvc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/quote"
)

type quoteClientTest func(ctx context.Context) (quote.Quote, error)

func (q quoteClientTest) RandomQuote(ctx context.Context) (quote.Quote, error) {
	return q(ctx)
}

var errStub = errors.New("stub")

func TestQuoteService(t *testing.T) {
	t.Run("GetRandomQuote", func(t *testing.T) {
		tests := map[string]struct {
			stubFunc func(ctx context.Context) (quote.Quote, error)
			wantErr  bool
		}{
			"ok": {
				stubFunc: func(ctx context.Context) (quote.Quote, error) {
					return quote.Quote{Author: "Ada", Text: "Keep going"}, nil
				},
				wantErr: false,
			},
			"error": {
				stubFunc: func(ctx context.Context) (quote.Quote, error) {
					return quote.Quote{}, errStub
				},
				wantErr: true,
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				svc := New(quoteClientTest(tc.stubFunc))

				_, err := svc.GetRandomQuote(context.Background())
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

package quotehdl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang/tutorial/todo/internal/apperr"
)

type quoteClientTest func(ctx context.Context) (string, error)

func (q quoteClientTest) GetRandomQuote(ctx context.Context) (string, error) {
	return q(ctx)
}

func TestQuoteHandler(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Parallel()

		cases := map[string]struct {
			stubFunc func(ctx context.Context) (string, error)
			wantErr  bool
		}{
			"ok": {
				stubFunc: func(ctx context.Context) (string, error) {
					return "Ada: Keep going", nil
				},
				wantErr: false,
			},
			"error": {
				stubFunc: func(ctx context.Context) (string, error) {
					return "", apperr.E(apperr.CodeUnknown, "Failed to get random quote", nil)
				},
				wantErr: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/quote", nil)
				h := New(quoteClientTest(tc.stubFunc))
				h.get(w, r)

				if !tc.wantErr {
					assert.Equal(t, http.StatusOK, w.Code)
					return
				}

				assert.Equal(t, http.StatusInternalServerError, w.Code)
			})
		}
	})
}

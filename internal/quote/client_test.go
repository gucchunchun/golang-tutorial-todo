package quote

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testQuoteText   = "Stay hungry, stay foolish."
	testQuoteAuthor = "Steve Jobs"
)

type table struct {
	name    string
	handler http.HandlerFunc
	want    Quote
	wantErr bool
}

func TestQuoteClient_RandomQuote(t *testing.T) {
	tests := []table{
		{
			name: "ok",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// if your client calls "/random", optionally assert:
				// if r.URL.Path != "/random" { t.Fatalf("path = %s", r.URL.Path) }
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode([]Quote{{
					Text:   testQuoteText,
					Author: testQuoteAuthor,
				}})
			},
			want: Quote{
				Text:   testQuoteText,
				Author: testQuoteAuthor,
			},
			wantErr: false,
		},
		{
			name: "bad status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "error", http.StatusInternalServerError)
			},
			want:    Quote{},
			wantErr: true, // should be true: non-2xx should return error
		},
		{
			name: "invalid json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{ "text": 123, `)) // malformed JSON
			},
			want:    Quote{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			// IMPORTANT: create the client with the test server base URL
			c := New(ts.URL, 10*time.Second)

			ctx := context.Background()
			got, err := c.RandomQuote(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Text, got.Text)
			assert.Equal(t, tt.want.Author, got.Author)
		})
	}
}

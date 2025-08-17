package quotes_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang/tutorial/todo/internal/quotes"
)

const (
	testQuoteText   = "Stay hungry, stay foolish."
	testQuoteAuthor = "Steve Jobs"
)

func TestHTTPClient_RandomQuote_OK(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("URL: %v\n", r.URL)
		if r.URL.Path != "/random" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]quotes.Quote{{
			Text:   testQuoteText,
			Author: testQuoteAuthor,
		}})
	}))
	defer ts.Close()

	c := quotes.NewHTTPClient(ts.URL, 10*time.Second)

	got, err := c.RandomQuote(context.Background())
	if err != nil {
		t.Fatalf("RandomQuote returned error: %v", err)
	}

	if got.Text != testQuoteText {
		t.Errorf("Text = %q; want %q", got.Text, testQuoteText)
	}
	if got.Author != testQuoteAuthor {
		t.Errorf("Author = %q; want %q", got.Author, testQuoteAuthor)
	}
}

func TestHTTPClient_RandomQuote_BadStatus(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := quotes.NewHTTPClient(ts.URL, 10*time.Second)

	_, err := c.RandomQuote(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHTTPClient_RandomQuote_InvalidJSON(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "text": 123, `))
	}))
	defer ts.Close()

	c := quotes.NewHTTPClient(ts.URL, 10*time.Second)

	_, err := c.RandomQuote(context.Background())
	if err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}
}

func TestHTTPClient_RandomQuote_Timeout(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(11 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"text": "Never gets here",
		})
	}))
	defer ts.Close()

	c := quotes.NewHTTPClient(ts.URL, 10*time.Millisecond)

	_, err := c.RandomQuote(context.Background())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

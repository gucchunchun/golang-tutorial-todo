package httpx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew_AndOptions(t *testing.T) {
	c, err := New("https://api.example.com",
		WithTimeout(3*time.Second),
		WithHeader("X-Foo", "Bar"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, c.HTTP)
	assert.Equal(t, c.HTTP.Timeout, 3*time.Second)
	assert.Equal(t, c.Headers.Get("X-Foo"), "Bar")

	// WithTransport
	rt := http.DefaultTransport
	c2, err := New("https://api.example.com", WithTransport(rt))
	assert.NoError(t, err)
	assert.Equal(t, c2.HTTP.Transport, rt)
}

func TestBuildURL(t *testing.T) {
	base, _ := url.Parse("https://example.com/api/")
	c := &Client{BaseURL: base}

	cases := []struct {
		rel  string
		want string
	}{
		{"/v1/random", "https://example.com/api/v1/random"},
		{"v1/random", "https://example.com/api/v1/random"},
		{"/v1//random/", "https://example.com/api/v1/random"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, c.buildURL(tc.rel))
	}
}

func TestDoJSON_Success_GET(t *testing.T) {
	type payload struct {
		Msg string `json:"msg"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		assert.Equal(t, r.URL.Path, "/hello")
		assert.NotZero(t, r.Header.Get("User-Agent"))
		assert.Equal(t, r.Method, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload{Msg: "ok"})
	}))
	defer ts.Close()

	c, err := New(ts.URL)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	var out payload
	err = c.DoJSON(context.Background(), http.MethodGet, "/hello", nil, &out)
	assert.NoError(t, err)
	assert.Equal(t, out.Msg, "ok")
}

func TestDoJSON_Success_POST_WithBody_AndHeaders(t *testing.T) {
	type in struct {
		Name string `json:"name"`
	}
	type out struct {
		Echo string `json:"echo"`
	}

	var sawContentType, sawCustomHeader bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if strings.HasPrefix(ct, "application/json") {
			sawContentType = true
		}
		if r.Header.Get("X-Token") == "abc123" {
			sawCustomHeader = true
		}
		// read body JSON and echo
		var got in
		bs, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bs, &got)
		_ = json.NewEncoder(w).Encode(out{Echo: got.Name})
	}))
	defer ts.Close()

	c, err := New(ts.URL, WithHeader("X-Token", "abc123"))
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	var res out
	err = c.DoJSON(context.Background(), http.MethodPost, "/echo", in{Name: "Yuna"}, &res)
	if err != nil {
		t.Fatalf("DoJSON error: %v", err)
	}
	if !sawContentType {
		t.Fatalf("expected Content-Type application/json to be set when body is present")
	}
	if !sawCustomHeader {
		t.Fatalf("expected custom header from client to be included")
	}
	if res.Echo != "Yuna" {
		t.Fatalf("echo = %q; want %q", res.Echo, "Yuna")
	}
}

func TestDoJSON_OutNil_DiscardsBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c, err := New(ts.URL)
	assert.NoError(t, err)

	err = c.DoJSON(context.Background(), http.MethodDelete, "/resource/1", nil, nil)
	assert.NoError(t, err)
}

func TestDoJSON_Non2xx_ResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusTeapot)
	}))
	defer ts.Close()

	c, err := New(ts.URL)
	assert.NoError(t, err)

	var dest struct{}
	err = c.DoJSON(context.Background(), http.MethodGet, "/bad", nil, &dest)
	assert.Error(t, err)

	var re *ResponseError
	assert.ErrorAs(t, err, &re)
	assert.Equal(t, re.StatusCode, http.StatusTeapot)
	assert.NotZero(t, re.Body)
}

package quotes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type HTTPClient struct {
	BaseURL   string
	HTTP      *http.Client
	UserAgent string
	Timeout   time.Duration
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: timeout,
		},
		UserAgent: "todo-cli/1.0",
		Timeout:   timeout,
	}
}

func (e *ResponseError) Error() string {
	return http.StatusText(e.StatusCode)
}

func (c *HTTPClient) RandomQuote(ctx context.Context) (Quote, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/random", nil)
	if err != nil {
		return Quote{}, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Quote{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Quote{}, &ResponseError{StatusCode: resp.StatusCode}
	}

	var quote []Quote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return Quote{}, err
	}

	return quote[0], nil
}

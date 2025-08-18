package quote

import (
	"context"
	"errors"
	"net/http"
	"time"

	"golang/tutorial/todo/internal/httpx"
)

type Client struct {
	BaseURL string
	HTTP    *httpx.Client
}

/*
Reference: O'REILLY「実用GO言語」11.2 p.258 RoundTripper
*/
type customRoundTripper struct {
	base http.RoundTripper
}

func (c customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// リクエスト前に実施したい処理
	resp, err := c.base.RoundTrip(req)
	// リクエスト後に実施したい処理
	return resp, err
}

func New(baseURL string, timeout time.Duration) *Client {
	/*
		Reference: O'REILLY「実用GO言語」11.1.3 p.257 http.Clientを作成してリクエスト
	*/
	client, err := httpx.New(baseURL, httpx.WithTimeout(timeout), httpx.WithTransport(customRoundTripper{http.DefaultTransport}))
	if err != nil {
		return nil
	}
	return &Client{
		BaseURL: baseURL,
		HTTP:    client,
	}
}

func (c *Client) RandomQuote(ctx context.Context) (Quote, error) {
	quotes := []Quote{}
	// HTTPリクエストの発行
	err := c.HTTP.DoJSON(ctx, http.MethodGet, "/random", nil, &quotes)
	if err != nil {
		return Quote{}, err
	}

	if len(quotes) == 0 {
		return Quote{}, errors.New("no quotes returned")
	}

	return quotes[0], nil
}

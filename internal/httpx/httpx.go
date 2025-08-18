package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Client struct {
	BaseURL *url.URL
	HTTP    *http.Client
	Headers http.Header
}

type Option func(*Client) error

func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if c.HTTP == nil {
			c.HTTP = &http.Client{Timeout: d}
		} else {
			c.HTTP.Timeout = d
		}
		return nil
	}
}

func WithHeader(k, v string) Option {
	return func(c *Client) error {
		c.Headers.Add(k, v)
		return nil
	}
}

func WithTransport(rt http.RoundTripper) Option {
	return func(c *Client) error {
		if c.HTTP == nil {
			c.HTTP = &http.Client{Transport: rt}
		} else {
			c.HTTP.Transport = rt
		}
		return nil
	}
}

func New(base string, opts ...Option) (*Client, error) {
	if !strings.HasPrefix(base, "http") {
		return nil, fmt.Errorf("base URL must include scheme, got %q", base)
	}
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	c := &Client{
		BaseURL: u,
		HTTP:    &http.Client{},
		Headers: make(http.Header),
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type ResponseError struct {
	StatusCode int
	Body       string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("httpx: unexpected status %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// パスを結合する
func (c *Client) buildURL(p string) string {
	u := *c.BaseURL
	u.Path = path.Join(strings.TrimRight(c.BaseURL.Path, "/"), strings.TrimLeft(p, "/"))
	return u.String()
}

func (c *Client) DoJSON(ctx context.Context, method, relPath string, body any, out any) error {
	// bodyをJSONに変換
	var rdr io.Reader
	if body != nil {
		pr, pw := io.Pipe()
		go func() {
			enc := json.NewEncoder(pw)
			pw.CloseWithError(enc.Encode(body))
		}()
		rdr = pr
	}

	/*
		Reference: O'REILLY「実用GO言語」11.1.1 p.255 http.Get
		http.GET contextなしでの呼び出し
		contextを持たないので途中でキャンセルできない
	*/
	// resp, err := http.Get(c.BaseURL+"/random")
	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(relPath), rdr)
	if err != nil {
		return err
	}

	/*
		Reference: O'REILLY「実用GO言語」11.1.1 p.255 http.Get
		http.GET contextなしでの呼び出し
		contextを持たないので途中でキャンセルできない
	*/
	for k, vv := range c.Headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "todo-cli/1.0")
	}

	/*
		Reference: O'REILLY「実用GO言語」11.1.1 p.257 http.Get
		クライアントとサーバー間のTCPコネクションをクローズしないとfile descriptorがリークする
	*/
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// デバッグのために読み込む
		slurp, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return &ResponseError{StatusCode: resp.StatusCode, Body: string(slurp)}
	}

	// レスポンスがない場合
	if out == nil {
		return nil
	}

	// レスポンスをJSONに変換
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}
	return nil
}

/*
コンパイル時にインターフェイスを満たしているかチェックすることができる
*/
type ClientInterface interface {
	DoJSON(ctx context.Context, method, relPath string, body any, out any) error
}

var _ ClientInterface = (*Client)(nil)

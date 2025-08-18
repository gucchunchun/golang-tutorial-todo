package limitmw

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newReq(path, ip string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("X-Forwarded-For", ip)
	return req
}

func newaHandler() http.Handler {
	return Limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func TestLimiter(t *testing.T) {
	ResetForTest(time.Hour, 2, time.Hour, 1)

	h := newaHandler()

	StartLimiterGC()

	t.Run("ip-limit", func(t *testing.T) {
		ip := "198.51.100.7"

		tests := map[string]struct {
			ip   string
			want int
		}{
			"ok":           {ip, http.StatusOK},
			"hit-ip-limit": {ip, http.StatusTooManyRequests},
		}

		for _, test := range tests {
			t.Run(test.ip, func(t *testing.T) {
				rec := httptest.NewRecorder()
				req := newReq("/", test.ip)
				h.ServeHTTP(rec, req)
				require.Equal(t, test.want, rec.Code)
			})
		}
	})
}

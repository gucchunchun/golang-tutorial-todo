package logmw_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"golang/tutorial/todo/internal/api/middleware/logmw"
)

func TestLog_EmitsStatusCodeAndLength(t *testing.T) {
	var buf bytes.Buffer
	origOut := log.Writer()
	origFlags := log.Flags()
	defer func() {
		log.SetOutput(origOut)
		log.SetFlags(origFlags)
	}()
	log.SetOutput(&buf)
	log.SetFlags(0)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Millisecond)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/log", nil)

	logmw.Logger(h).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
	assert.Equal(t, "hello world", rec.Body.String())

	out := buf.String()
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "GET")
	assert.Contains(t, out, "/log")
	assert.Contains(t, out, "200")
	assert.Contains(t, out, "11B")

	hasDurationSuffix := strings.Contains(out, "ns") ||
		strings.Contains(out, "µs") ||
		strings.Contains(out, "ms") ||
		strings.Contains(out, "s")
	assert.True(t, hasDurationSuffix, "expected a duration in the log output, got: %q", out)
}

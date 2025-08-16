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

func TestLog_PassThroughResponse(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("hello"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	logmw.Log(h).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.Equal(t, "ok", rec.Header().Get("X-Test"))
	assert.Equal(t, "hello", rec.Body.String())
}

func TestLog_EmitsMethodPathAndDuration(t *testing.T) {
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
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/log", nil)

	logmw.Log(h).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	out := buf.String()
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "POST")
	assert.Contains(t, out, "/log")

	hasDurationSuffix := strings.Contains(out, "ns") ||
		strings.Contains(out, "µs") ||
		strings.Contains(out, "ms") ||
		strings.Contains(out, "s")
	assert.True(t, hasDurationSuffix, "expected a duration in the log output, got: %q", out)
}

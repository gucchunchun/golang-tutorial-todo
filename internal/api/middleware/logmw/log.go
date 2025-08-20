package logmw

import (
	"net/http"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

// Reference: O'REILLY「実用GO言語」10.5.2 p.244
// Reference: O'REILLY「実用GO言語」12.4 p.281
type loggingResponseWriter struct {
	statusCode int
	length     int
	writer     http.ResponseWriter
	request    *http.Request
	start      time.Time
}

func newLoggingResponseWriter(w http.ResponseWriter, r *http.Request) *loggingResponseWriter {
	return &loggingResponseWriter{writer: w, request: r, start: time.Now()}
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.writer.Header()
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.writer.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	} else if lrw.statusCode >= 400 {
		log.Printf("client error: %s", b)
	}
	size, err := lrw.writer.Write(b)
	lrw.length += size
	return size, err
}

func (lrw *loggingResponseWriter) MarshalZerologObject(e *zerolog.Event) {
	e.Str("method", lrw.request.Method)
	e.Str("path", lrw.request.URL.Path)
	e.Int64("request_size", lrw.request.ContentLength)
	e.Int("status_code", lrw.statusCode)
	e.Int("response_size", lrw.length)
	e.Str("referer", lrw.request.Referer())
	e.Str("latency", time.Since(lrw.start).String())
}

type contextKey string

const logKey contextKey = "log"

func Logger(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get("X-Request-Id")
			if traceID == "" {
				// ユニークIDを生成
				traceID = xid.New().String()
			}
			logger := logger.With().Str("trace_id", traceID).Logger()
			ctx := context.WithValue(r.Context(), logKey, &logger)

			writer := newLoggingResponseWriter(w, r)
			next.ServeHTTP(writer, r.WithContext(ctx))
			logger.Info().Object("httpRequest", writer).Send()
		})
	}
}

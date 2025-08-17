package logmw

import (
	"log"
	"net/http"
	"time"
)

// Reference: O'REILLY「実用GO言語」10.5.2 p.244
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{ResponseWriter: w}
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	} else if lrw.statusCode >= 400 {
		log.Printf("client error: %s", b)
	}
	size, err := lrw.ResponseWriter.Write(b)
	lrw.length += size
	return size, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := NewLoggingResponseWriter(w)

		next.ServeHTTP(sw, r)

		log.Printf(
			"%s %s %d %dB %s",
			r.Method,
			r.URL.Path,
			sw.statusCode,
			sw.length,
			time.Since(start),
		)
	})
}

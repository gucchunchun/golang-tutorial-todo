package logmw

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	// NOTE: http.ResponseWriterを埋め込む = extendsのような感じ
	//       メゾッドをオーバーライドすることが可能
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.length += size
	return size, err
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r)

		log.Println("logout")

		log.Printf(
			"%s %s %d %dB %s",
			r.Method,
			r.URL.Path,
			sw.status,
			sw.length,
			time.Since(start),
		)
	})
}

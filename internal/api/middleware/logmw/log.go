package logmw

import (
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

// Reference: O'REILLY「実用GO言語」10.5.2 p.244
// Reference: O'REILLY「実用GO言語」12.4 p.281
type loggingResponseWriter struct {
	statusCode int
	length     int
	http.ResponseWriter
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(b)
	lrw.length += n
	return n, err
}

type contextKey string

const logKey contextKey = "log"

var (
	accessFile     *os.File
	accessFileOnce sync.Once
	accessFileErr  error
	writeMu        sync.Mutex
)

func openAccessFile() (*os.File, error) {
	accessFileOnce.Do(func() {
		path := os.Getenv("ACCESS_LOG_FILE")
		if path == "" {
			// choose container path if exists
			if st, err := os.Stat("/var/log/app"); err == nil && st.IsDir() {
				path = "/var/log/app/access.jsonl"
			} else {
				path = "./logs/access.jsonl"
			}
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			accessFileErr = err
			return
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			accessFileErr = err
			return
		}
		accessFile = f
	})
	return accessFile, accessFileErr
}

// フォーマット:
// {"ts":"<RFC3339>","method":"GET","path":"/tasks/123","route":"/tasks/{id}","status":200,"bytes":412,"latency_ms":23,"ip":"203.0.113.5","ua":"curl/8.6.0","trace":"xid_abc123"}
func Logger(base *zerolog.Logger) func(next http.Handler) http.Handler {
	// アクセスログ出力用のロガーを作成
	var mw io.Writer = os.Stdout
	if f, err := openAccessFile(); err == nil && f != nil {
		mw = io.MultiWriter(os.Stdout, f)
	} else if err != nil && !errors.Is(err, os.ErrPermission) {
		base.Error().Err(err).Msg("failed to open access log file")
	}
	accessLogger := zerolog.New(mw)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// traceID ：一連の関連する情報に対して付与される識別子
			traceID := r.Header.Get("X-Request-Id")
			if traceID == "" {
				traceID = xid.New().String()
			}

			reqLogger := base.With().Str("trace", traceID).Logger()
			ctx := context.WithValue(r.Context(), logKey, &reqLogger)

			lrw := newLoggingResponseWriter(w)
			next.ServeHTTP(lrw, r.WithContext(ctx))

			if lrw.statusCode == 0 {
				lrw.statusCode = http.StatusOK
			}
			path := r.URL.Path
			route := routeFromContext(r.Context())
			if route == "" {
				route = path
			}
			entry := accessLogEntry{
				TS:        start.UTC(),
				Method:    r.Method,
				Path:      path,
				Route:     route,
				Status:    lrw.statusCode,
				Bytes:     lrw.length,
				LatencyMS: time.Since(start).Milliseconds(),
				IP:        clientIP(r),
				UA:        r.UserAgent(),
				Trace:     traceID,
			}
			accessLogger.Info().EmbedObject(entry).Send()
		})
	}
}

type accessLogEntry struct {
	TS        time.Time
	Method    string
	Path      string
	Route     string
	Status    int
	Bytes     int
	LatencyMS int64
	IP        string
	UA        string
	Trace     string
}

func (a accessLogEntry) MarshalZerologObject(e *zerolog.Event) {
	e.Str("ts", a.TS.Format(time.RFC3339Nano)).
		Str("method", a.Method).
		Str("path", a.Path).
		Str("route", a.Route).
		Int("status", a.Status).
		Int("bytes", a.Bytes).
		Int64("latency_ms", a.LatencyMS).
		Str("ip", a.IP).
		Str("ua", a.UA).
		Str("trace", a.Trace)
}

// routeFromContext コンテキスト内にあるルートパターンの取得をトライ
func routeFromContext(ctx context.Context) string {
	if v := ctx.Value(contextKey("route_pattern")); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// clientIP X-Forwarded-For または RemoteAddr からIPアドレスを取得する
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func Flush() error {
	if accessFile == nil {
		return nil
	}
	writeMu.Lock()
	defer writeMu.Unlock()
	return accessFile.Sync()
}

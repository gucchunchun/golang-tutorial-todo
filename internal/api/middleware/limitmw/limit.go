package limitmw

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Reference: O'REILLY「実用GO言語」10.562 p.248 レートリミット
// 設定値
const (
	// 全体: 1,000 req/sec, バースト 200
	globalEvery = time.Millisecond // 1msごとに1トークン ≒ 1000 rps
	globalBurst = 200

	// ユーザー: 1秒に1回、バースト1（= 実質きっちり1秒に1回）
	userEvery = time.Second
	userBurst = 1

	// この待ち時間を超えるなら 429 を返す（待たせない）
	maxWait = 0 * time.Millisecond

	// ユーザーエントリのGarbageCollection（この時間アクセスなければ破棄）
	// Time To Live
	visitorTTL = 10 * time.Minute
)

var globalLimiter = rate.NewLimiter(rate.Every(globalEvery), globalBurst)

var newUserLimiter = func() *rate.Limiter { return rate.NewLimiter(rate.Every(userEvery), userBurst) }

type visitorLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type visitorsStore struct {
	mu       sync.Mutex
	visitors map[string]*visitorLimiter
}

var store = &visitorsStore{visitors: make(map[string]*visitorLimiter)}

func (s *visitorsStore) get(key string) *rate.Limiter {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	vl, ok := s.visitors[key]
	if !ok {
		vl = &visitorLimiter{
			limiter:  newUserLimiter(),
			lastSeen: now,
		}
		s.visitors[key] = vl
	} else {
		vl.lastSeen = now
	}
	return vl.limiter
}

func (s *visitorsStore) gc() {
	ticker := time.NewTicker(5 * time.Minute)
	for now := range ticker.C {
		s.mu.Lock()
		for k, v := range s.visitors {
			if now.Sub(v.lastSeen) > visitorTTL {
				delete(s.visitors, k)
			}
		}
		s.mu.Unlock()
	}
}

// ミドルウェア呼び出し時に起動する
func StartLimiterGC() {
	go store.gc()
}

func Limiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1) グローバル制限
		if !allowOrAdviseRetry(w, globalLimiter) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		// 2) ユーザー（IP）ごとの制限
		key := clientKey(r) // IPベース。認証があるならユーザーIDに置換推奨
		if !allowOrAdviseRetry(w, store.get(key)) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Reserve() でDelayを求め、待たせずに 429 + Retry-After を返す。
// 受け入れる場合は true、拒否してレスポンスを書いた場合は false を返す。
func allowOrAdviseRetry(w http.ResponseWriter, lim *rate.Limiter) bool {
	res := lim.Reserve() // ここでトークンを仮確保
	if !res.OK() {
		// 物理的に受け付けられない（上限が小さすぎる等）
		w.Header().Set("Retry-After", "1")
		return false
	}
	delay := res.Delay()
	if delay <= maxWait {
		// すぐ通せるので予約を有効化してOK
		return true
	}
	// 通さないので予約をキャンセル（トークンを戻す）
	res.Cancel()

	// RFC 7231 準拠の秒指定（四捨五入より安全に切り上げ）
	sec := max(int(math.Ceil(delay.Seconds())), 1)
	w.Header().Set("Retry-After", strconv.Itoa(sec))
	return false
}

// プロキシ下対応（X-Forwarded-For → 最初の値）: 必要に応じて調整
func clientKey(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

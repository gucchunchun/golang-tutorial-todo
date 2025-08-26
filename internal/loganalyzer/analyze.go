package loganalyzer

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"time"

	"gonum.org/v1/gonum/stat"
)

// Stats 集計結果
type Stats struct {
	Count    int
	From, To time.Time
	ByStatus map[int]int
	ByClass  map[string]int // 2xx/3xx/4xx/5xx/other
	ByGroup  map[string]int // route or path
	LatMin   int64
	LatAvg   float64
	LatP50   float64
	LatP95   float64
	LatP99   float64
	LatMax   int64
}

// Options 解析オプション
type Options struct {
	GroupBy        string // route | path
	TopN           int
	EstimatedLines int // 推定行数 (0 の場合はデフォルト 4096 キャパ)
}

// KV 上位エンドポイント表示用
type KV struct {
	Key   string
	Count int
}

// accessLogRaw fallback 用構造体 (encoding/json 使用時)
type accessLogRaw struct {
	TS        string `json:"ts"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Route     string `json:"route"`
	Status    int    `json:"status"`
	Bytes     int    `json:"bytes"`
	LatencyMS int64  `json:"latency_ms"`
}

// interner 重複文字列を一度だけ保持
type interner struct{ m map[string]string }

func newInterner(capHint int) *interner { return &interner{m: make(map[string]string, capHint)} }
func (in *interner) intern(s string) string {
	if v, ok := in.m[s]; ok {
		return v
	}
	in.m[s] = s
	return s
}

// Analyze ログをストリーミング解析
func Analyze(r io.Reader, opt Options) (Stats, []KV, error) {
	sc := bufio.NewScanner(r)
	s := Stats{ByStatus: make(map[int]int, 16), ByClass: make(map[string]int, 8), ByGroup: make(map[string]int, 256)}
	var firstTS time.Time
	var firstParsed bool
	var lastTSStr string
	routeIntern := newInterner(512)
	pathIntern := newInterner(512)
	// 使い回す Unmarshal 用構造体
	var raw accessLogRaw
	var (
		latSum  int64
		capHint int
	)
	if opt.EstimatedLines > 0 {
		capHint = opt.EstimatedLines
	} else {
		capHint = 4096
	}
	latencies := make([]float64, 0, capHint)
	for sc.Scan() {
		b := sc.Bytes()
		raw = accessLogRaw{} // 前回値を消す
		if err := json.Unmarshal(b, &raw); err != nil {
			continue
		}
		if raw.TS == "" {
			continue
		}
		if !firstParsed {
			ts, err := time.Parse(time.RFC3339Nano, raw.TS)
			if err != nil {
				continue
			}
			firstTS = ts
			s.From = ts
			firstParsed = true
		}
		lastTSStr = raw.TS
		s.Count++
		s.ByStatus[raw.Status]++
		s.ByClass[classOf(raw.Status)]++
		rStr := routeIntern.intern(raw.Route)
		pStr := pathIntern.intern(raw.Path)
		key := rStr
		if opt.GroupBy == "path" {
			key = pStr
		}
		s.ByGroup[key]++
		latSum += raw.LatencyMS
		latencies = append(latencies, float64(raw.LatencyMS))
	}
	if err := sc.Err(); err != nil && !errors.Is(err, io.EOF) {
		return Stats{}, nil, err
	}
	if s.Count == 0 {
		return s, nil, nil
	}
	if lastTSStr != "" {
		if ts, err := time.Parse(time.RFC3339Nano, lastTSStr); err == nil {
			s.To = ts
		} else {
			s.To = firstTS
		}
	} else {
		s.To = firstTS
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	s.LatMin = int64(latencies[0])
	s.LatMax = int64(latencies[len(latencies)-1])
	s.LatAvg = float64(latSum) / float64(len(latencies))
	s.LatP50 = stat.Quantile(0.50, stat.Empirical, latencies, nil)
	s.LatP95 = stat.Quantile(0.95, stat.Empirical, latencies, nil)
	s.LatP99 = stat.Quantile(0.99, stat.Empirical, latencies, nil)
	kvs := make([]KV, 0, len(s.ByGroup))
	for k, c := range s.ByGroup {
		kvs = append(kvs, KV{Key: k, Count: c})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].Count > kvs[j].Count })
	if opt.TopN > 0 && len(kvs) > opt.TopN {
		kvs = kvs[:opt.TopN]
	}
	return s, kvs, nil
}

func classOf(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500 && status < 600:
		return "5xx"
	default:
		return "other"
	}
}

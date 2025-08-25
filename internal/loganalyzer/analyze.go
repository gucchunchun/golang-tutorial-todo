package loganalyzer

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"time"
)

type AccessLog struct {
	TS     time.Time `json:"ts"`
	Method string    `json:"method"`
	Path   string    `json:"path"`
	Route  string    `json:"route"`
	Status int       `json:"status"`
	Bytes  int       `json:"bytes"`
	// LatencyMS time.Duration の Milliseconds メゾットが int64 をリターンするため
	LatencyMS int64 `json:"latency_ms"`
}

type Stats struct {
	Count    int
	From, To time.Time
	ByStatus map[int]int
	ByClass  map[string]int // "2xx","4xx","5xx" 等
	ByGroup  map[string]int // route or path
	LatMin   int64
	LatAvg   float64
	LatP50   int64
	LatP95   int64
	LatP99   int64
	LatMax   int64
}

type Options struct {
	GroupBy string // "route" or "path"
	TopN    int
}

type KV struct {
	Key   string
	Count int
}

func Analyze(r io.Reader, opt Options) (Stats, []KV, error) {
	sc := bufio.NewScanner(r)
	var (
		latencies []int64
		s         = Stats{
			ByStatus: make(map[int]int),
			ByClass:  make(map[string]int),
			ByGroup:  make(map[string]int),
		}
		firstTS time.Time
	)

	/*
		sc.Scan() 読み込みトークンが残っている限り、Bytes()/Text() で次のトークンを読み込み可能。
		読み込めるトークンがなくなった場合にfalseを返す。
	*/
	for sc.Scan() {
		var al AccessLog
		// 壊れた行はスキップ
		if err := json.Unmarshal(sc.Bytes(), &al); err != nil {
			continue
		}
		if al.TS.IsZero() {
			continue
		}
		if firstTS.IsZero() {
			firstTS = al.TS
			s.From = al.TS
		}
		s.To = al.TS
		s.Count++
		s.ByStatus[al.Status]++
		s.ByClass[classOf(al.Status)]++
		key := al.Route
		if opt.GroupBy == "path" {
			key = al.Path
		}
		s.ByGroup[key]++
		latencies = append(latencies, al.LatencyMS)
	}
	if err := sc.Err(); err != nil && !errors.Is(err, io.EOF) {
		return Stats{}, nil, err
	}
	if s.Count == 0 {
		return s, nil, nil
	}

	// レイテンシーの最小・最大値
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	s.LatMin = latencies[0]
	s.LatMax = latencies[len(latencies)-1]

	// 平均
	var sum int64
	for _, v := range latencies {
		sum += v
	}
	s.LatAvg = float64(sum) / float64(len(latencies))

	// パーセンタイル
	p := func(q float64) int64 {
		n := len(latencies)
		if n == 0 {
			return 0
		}
		if q <= 0 {
			return latencies[0]
		}
		if q >= 1 {
			return latencies[n-1]
		}
		idx := q * float64(n-1)
		lowerIdx := int(idx)
		upperIdx := lowerIdx + 1
		if upperIdx >= n {
			upperIdx = lowerIdx
		}
		fraction := idx - float64(lowerIdx)

		lowerValue := float64(latencies[lowerIdx])
		upperValue := float64(latencies[upperIdx])

		return int64(lowerValue + (upperValue-lowerValue)*fraction)
	}
	s.LatP50 = p(0.50)
	s.LatP95 = p(0.95)
	s.LatP99 = p(0.99)

	// 上位 N（ByGroup）
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

// classOf ステータスごとにクラスを返す
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

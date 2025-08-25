package loganalyzer

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze_OK(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		opts     Options
		expected Stats
	}
	tests := []testcase{
		{
			name: "two valid lines same route",
			input: `{"ts":"2025-08-24T23:00:00.000Z","method":"GET","path":"/tasks","route":"/tasks","status":200,"bytes":100,"latency_ms":12}
							{"ts":"2025-08-24T23:00:00.500Z","method":"POST","path":"/tasks","route":"/tasks","status":201,"bytes":150,"latency_ms":12}
							`,
			opts: Options{GroupBy: "route", TopN: 5},
			expected: Stats{
				Count: 2,
				From:  time.Date(2025, 8, 24, 23, 0, 0, 0, time.UTC),
				To:    time.Date(2025, 8, 24, 23, 0, 0, 500_000_000, time.UTC),
				ByStatus: map[int]int{
					200: 1,
					201: 1,
				},
				ByClass: map[string]int{
					"2xx": 2,
				},
				ByGroup: map[string]int{
					"/tasks": 2,
				},
				LatMin: 12,
				LatAvg: 12.0,
				LatP50: 12,
				LatP95: 12,
				LatP99: 12,
				LatMax: 12,
			},
		},
		{
			name: "single valid line",
			input: `{"ts":"2025-08-25T00:10:00.000Z","method":"GET","path":"/tasks","route":"/tasks","status":200,"bytes":42,"latency_ms":7}
`,
			opts: Options{GroupBy: "route", TopN: 5},
			expected: Stats{
				Count: 1,
				From:  time.Date(2025, 8, 25, 0, 10, 0, 0, time.UTC),
				To:    time.Date(2025, 8, 25, 0, 10, 0, 0, time.UTC),
				ByStatus: map[int]int{
					200: 1,
				},
				ByClass: map[string]int{
					"2xx": 1,
				},
				ByGroup: map[string]int{
					"/tasks": 1,
				},
				LatMin: 7,
				LatAvg: 7.0,
				LatP50: 7,
				LatP95: 7,
				LatP99: 7,
				LatMax: 7,
			},
		},
		{
			name: "valid log lines",
			input: `{"ts":"2025-08-24T23:00:00.000Z","method":"GET","path":"/tasks","route":"/tasks","status":200,"bytes":100,"latency_ms":12}
							THIS IS NOT JSON
							{"ts":"","method":"GET","path":"/broken","route":"/broken","status":200,"bytes":0,"latency_ms":1}
							`,
			opts: Options{GroupBy: "route", TopN: 5},
			expected: Stats{
				Count: 1,
				From:  time.Date(2025, 8, 24, 23, 0, 0, 0, time.UTC),
				To:    time.Date(2025, 8, 24, 23, 0, 0, 0, time.UTC),
				ByStatus: map[int]int{
					200: 1,
				},
				ByClass: map[string]int{
					"2xx": 1,
				},
				ByGroup: map[string]int{
					"/tasks": 1,
				},
				LatMin: 12,
				LatAvg: 12.0,
				LatP50: 12,
				LatP95: 12,
				LatP99: 12,
				LatMax: 12,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := strings.NewReader(test.input)
			got, top, err := Analyze(in, test.opts)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, got)
			assert.LessOrEqual(t, len(top), test.opts.TopN)
		})
	}
}

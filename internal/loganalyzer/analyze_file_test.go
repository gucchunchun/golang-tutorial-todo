package loganalyzer

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyzeFile_OK 通常ファイルとgzipファイルの両方で正しく解析できることを確認
func TestAnalyzeFile_OK(t *testing.T) {
	dir := t.TempDir()
	lines := []string{
		`{"ts":"2025-08-24T23:00:00.000Z","method":"GET","path":"/tasks","route":"/tasks","status":200,"bytes":100,"latency_ms":12}`,
		`THIS IS BROKEN JSON`,
		`{"ts":"","method":"GET","path":"/skip","route":"/skip","status":200,"bytes":0,"latency_ms":1}`,
		`{"ts":"2025-08-24T23:00:00.500Z","method":"POST","path":"/tasks","route":"/tasks","status":201,"bytes":150,"latency_ms":10}`,
	}

	writePlain := func(path string) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		for _, l := range lines {
			if _, err := w.WriteString(l + "\n"); err != nil {
				return err
			}
		}
		return w.Flush()
	}
	writeGzip := func(path string) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		gw := gzip.NewWriter(f)
		for _, l := range lines {
			if _, err := gw.Write([]byte(l + "\n")); err != nil {
				return err
			}
		}
		return gw.Close()
	}

	plainPath := filepath.Join(dir, "access.jsonl")
	err := writePlain(plainPath)
	require.NoError(t, err)

	gzPath := filepath.Join(dir, "access.jsonl.gz")
	err = writeGzip(gzPath)
	require.NoError(t, err)

	expectFrom := time.Date(2025, 8, 24, 23, 0, 0, 0, time.UTC)
	expectTo := time.Date(2025, 8, 24, 23, 0, 0, 500_000_000, time.UTC)

	check := func(label, path string) {
		t.Run(label, func(t *testing.T) {
			got, top, err := AnalyzeFile(path, Options{GroupBy: "route", TopN: 10})
			assert.NoError(t, err)

			assert.Equal(t, 2, got.Count)
			assert.Equal(t, expectFrom, got.From)
			assert.Equal(t, expectTo, got.To)

			assert.Equal(t, 1, got.ByStatus[200])
			assert.Equal(t, 1, got.ByStatus[201])
			assert.Equal(t, 2, got.ByClass["2xx"])

			assert.Equal(t, 2, got.ByGroup["/tasks"])

			assert.Equal(t, int64(10), got.LatMin)
			assert.Equal(t, int64(12), got.LatMax)

			assert.NotEmpty(t, top)
			assert.LessOrEqual(t, len(top), 10)
		})
	}
	check("plain", plainPath)
	check("gzip", gzPath)
}

// BenchmarkAnalyzeFile_Large サイズの大きいログファイルを解析するベンチマーク
func BenchmarkAnalyzeFile_Large(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "large.jsonl")
	const N = 100_000

	check := func(label string, estimatedLines int) {
		b.Run(label, func(b *testing.B) {
			// テストファイルの作成
			f, err := os.Create(path)
			require.NoError(b, err)
			w := bufio.NewWriter(f)
			rnd := rand.New(rand.NewSource(42))
			routes := []string{"/tasks", "/tasks/{id}", "/health", "/quotes"}
			statuses := []int{200, 201, 404, 500}
			base := time.Date(2025, 8, 25, 0, 0, 0, 0, time.UTC)
			for i := 0; i < N; i++ {
				ts := base.Add(time.Duration(i) * time.Millisecond).UTC().Format(time.RFC3339Nano)
				route := routes[rnd.Intn(len(routes))]
				pathVal := route
				st := statuses[rnd.Intn(len(statuses))]
				lat := rnd.Intn(900) + 1 // 1..900 ms
				line := fmt.Sprintf(`{"ts":"%s","method":"GET","path":"%s","route":"%s","status":%d,"bytes":%d,"latency_ms":%d}`, ts, pathVal, route, st, rnd.Intn(2048), lat)
				_, err := w.WriteString(line + "\n")
				require.NoError(b, err)
			}
			require.NoError(b, w.Flush())
			require.NoError(b, f.Close())

			// 準備処理に時間がかかるため、ベンチマーク対象から除外
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _, err := AnalyzeFile(path, Options{GroupBy: "route", TopN: 5, EstimatedLines: estimatedLines})
				assert.NoError(b, err)
			}
		})
	}

	check("without estimation", 0)
	check("with exact estimation", N)
	check("with smaller estimation", N/2)
	check("with larger estimation", N*10)
}

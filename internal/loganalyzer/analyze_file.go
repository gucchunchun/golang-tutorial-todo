package loganalyzer

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func AnalyzeFile(path string, opt Options) (Stats, []KV, error) {
	r, err := openMaybeGzip(path)
	if err != nil {
		return Stats{}, nil, err
	}
	defer r.Close()
	return Analyze(r, opt)
}

// openMaybeGzip gzipファイルを含めて開くことが可能
func openMaybeGzip(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if ext := filepath.Ext(path); ext == ".gz" {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		return struct {
			io.Reader
			io.Closer
		}{Reader: gr, Closer: f}, nil
	}
	return f, nil
}

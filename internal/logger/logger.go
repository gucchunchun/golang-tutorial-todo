package logger

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type Config struct {
	Service         string
	FilePath        string
	Level           zerolog.Level
	ConsoleWarnOnly bool
	UseGlobal       bool
}

type warnOnlyLevelWriter struct{ w io.Writer }

func (w warnOnlyLevelWriter) Write(p []byte) (int, error) { return w.w.Write(p) }
func (w warnOnlyLevelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level >= zerolog.WarnLevel {
		return w.w.Write(p)
	}
	// pretend it was written to satisfy callers
	return len(p), nil
}

func Setup(cfg Config) (*zerolog.Logger, func(), error) {
	zerolog.SetGlobalLevel(cfg.Level)

	// デフォルトのフィールド名を変更できる
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "p"
	zerolog.MessageFieldName = "m"

	// Console
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	var consoleSink io.Writer = console
	if cfg.ConsoleWarnOnly {
		consoleSink = warnOnlyLevelWriter{w: console}
	}

	// File
	var file *os.File
	var writers []io.Writer
	writers = append(writers, consoleSink)

	if cfg.FilePath != "" {
		f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return &zerolog.Logger{}, nil, err
		}
		file = f
		writers = append(writers, f)
	}

	mw := zerolog.MultiLevelWriter(writers...)
	base := zerolog.New(mw).
		With().
		Timestamp().
		Str("service", cfg.Service).
		Logger()

	log.SetFlags(0)
	log.SetOutput(console)

	if cfg.UseGlobal {
		zlog.Logger = base
	}

	cleanup := func() {
		if file != nil {
			_ = file.Close()
		}
	}

	return &base, cleanup, nil
}

// contextにloggerを埋め込む
func IntoContext(ctx context.Context, lg zerolog.Logger) context.Context {
	return lg.WithContext(ctx)
}

// contextからloggerを取り出す
func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

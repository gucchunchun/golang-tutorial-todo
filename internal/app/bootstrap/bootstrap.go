package bootstrap

/*
インフラのセットアップをappではなく、別パッケージに分けることで、
- アプリケーションのオーケストレーションから分離できる
- 様々な場所で使用できるかつ、testが単体でできるため容易
- インフラ関連パッケージのimportを統一できる
*/

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"golang/tutorial/todo/internal/config"
	"golang/tutorial/todo/internal/db"
	"golang/tutorial/todo/internal/logger"
)

type Options struct {
	Service         string
	LogFile         string
	Level           zerolog.Level
	ConsoleWarnOnly bool
	UseGlobal       bool
}

type Deps struct {
	Ctx     context.Context
	Log     *zerolog.Logger
	DB      *db.Client
	Cleanup func()
}

func Init(ctx context.Context, opt Options) (*Deps, error) {
	// 環境変数の読み込み
	_ = godotenv.Load()

	// ロガーの初期化
	lg, cleanup, err := logger.Setup(logger.Config{
		Service:         opt.Service,
		FilePath:        opt.LogFile,
		Level:           opt.Level,
		ConsoleWarnOnly: opt.ConsoleWarnOnly,
		UseGlobal:       opt.UseGlobal,
	})
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		return nil, err
	}

	// DBの初期化
	cfg := config.NewDBConfig()
	dbc := &db.Client{}
	if err := dbc.Connect(cfg); err != nil {
		return nil, err
	}
	lg.Info().Msg("DB connected")

	// SIGINT: Ctrl + C, SIGTERM: Dockerなどコンテナからの終了シグナル
	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return &Deps{
		Ctx: sigCtx,
		Log: lg,
		DB:  dbc,
		Cleanup: func() {
			stop()
			dbc.Close()
			if cleanup != nil {
				cleanup()
			}
		},
	}, nil
}

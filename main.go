package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"golang/tutorial/todo/cmd"
	"golang/tutorial/todo/internal/logger"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// loggerの初期化
	lg, cleanup, err := logger.Setup(logger.Config{
		Service:         "todo-service",
		FilePath:        "./logs/myapp.log",
		Level:           zerolog.DebugLevel,
		ConsoleWarnOnly: false,
		UseGlobal:       true,
	})
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()
	if err != nil {
		fmt.Println("logger setup: %w", err)
		os.Exit(1)
	}

	// コブラコマンドの実行
	log.Default().Println("cobra executed")
	ctx := lg.WithContext(context.Background())
	if err := cmd.Execute(ctx); err != nil {
		fmt.Println(err)
		/*
			Reference: O'REILLY「実用GO言語」13.5 p.285 強制終了しても良い場所
			main(), init(), Must()など初期化を行う関数内ではpanic(), log.Fatal(), os.Exit()などを使い処理を終了させることができる。
		*/
		os.Exit(1)
	}
}

package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"

	"golang/tutorial/todo/cmd"
	"golang/tutorial/todo/internal/quotes"
	"golang/tutorial/todo/internal/services/task"
	"golang/tutorial/todo/internal/storage"
)

func main() {
	// 環境変数の読み込み
	godotenv.Load()

	// サービスの初期化
	quoteClient := quotes.NewHTTPClient(os.Getenv("QUOTES_BASE_URL"), 10*time.Second)
	storageClient := storage.New("tasks.json")
	taskService := task.NewService(quoteClient, storageClient)

	// コマンドの実行
	cmd.Execute(taskService)
}

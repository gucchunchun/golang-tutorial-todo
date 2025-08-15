package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"golang/tutorial/todo/cmd"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// コマンドの実行
	cmd.Execute()
}

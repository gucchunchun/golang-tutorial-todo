package main

import (
	"context"
	"log"

	"golang/tutorial/todo/internal/app/server"
)

/*
Reference: O'REILLY「実用GO言語」13.5 p.285 強制終了しても良い場所
main(), init(), Must()など初期化を行う関数内ではpanic(), log.Fatal(), os.Exit()などを使い処理を終了させることができる。
*/
func main() {
	ctx := context.Background()
	if err := server.RunHTTPServer(ctx); err != nil {
		log.Fatal(err)
	}
}

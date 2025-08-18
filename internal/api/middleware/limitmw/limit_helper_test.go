package limitmw

import (
	"time"

	"golang.org/x/time/rate"
)

/*
Reference: O'REILLY「実用GO言語」13.4 p.298 ヘルパー関数
テスト用の共通処理を定義するヘルパー関数は "helper_test.go" / "{package}_helper_test.go"
"**_test.go" にまとめることで他パッケージからはimportできない
*/

// limiter の内部状態をリセットし、グローバル/ユーザーの補充間隔＆バーストをテスト用に上書きする
func ResetForTest(globalEvery time.Duration, globalBurst int, userEvery time.Duration, userBurst int) {
	// グローバル limiter を置き換え
	globalLimiter = rate.NewLimiter(rate.Every(globalEvery), globalBurst)

	// ストアを新品に
	store = &visitorsStore{visitors: make(map[string]*visitorLimiter)}

	// ユーザー用 limiter の生成方法を差し替え
	newUserLimiter = func() *rate.Limiter {
		return rate.NewLimiter(rate.Every(userEvery), userBurst)
	}
}

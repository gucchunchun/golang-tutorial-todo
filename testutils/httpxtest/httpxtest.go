package httpxtest

import (
	"context"

	"golang/tutorial/todo/internal/httpx"
)

/*
Reference: O'REILLY「実用GO言語」13.4.2 p.299 テストヘルパー用のパッケージを作成する
**_test.goでないため、その他パッケージからimport可能 = 服うすうパッケージにまたがるヘルパー関数を作成可能。
それぞれのパッケージ内の**_test.goからのみimportされればビルド時には含まれないため問題にならない。
あるパッケージの機能をモックする場合は、{package}testという名称が一般的。
*/

// メゾットが一つの場合：
type TestClient func(ctx context.Context, method, relPath string, body any, out any) error

func (t TestClient) DoJSON(ctx context.Context, method, relPath string, body any, out any) error {
	return t(ctx, method, relPath, body, out)
}

// メゾットが複数の場合：
// type TestClient struct {
// 	DoJSONFunc func(ctx context.Context, method, relPath string, body any, out any) error
// }

// func (t *TestClient) DoJSON(ctx context.Context, method, relPath string, body any, out any) error {
// 	if t.DoJSONFunc == nil {
// 		panic("DoJSONFunc is not set")
// 	}
// 	return t.DoJSONFunc(ctx, method, relPath, body, out)
// }

// インターフェイスを満たしているかチェック
var _ httpx.ClientInterface = (*TestClient)(nil)

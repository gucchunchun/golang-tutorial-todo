package httpxtest

import (
	"context"
	"encoding/json"
	"fmt"
)

/*
Reference: O'REILLY「実用GO言語」13.12 p.323 ドキュメントに出力するExampleをコードに記述する
ExampleF: function
ExampleT: struct T
ExampleT_M: (struct T) method M
*/

func ExampleTestClient() {
	tc := TestClient(func(ctx context.Context, method, relPath string, body any, out any) error {
		// ここにテスト用のコードを書く
		return nil
	})
	fmt.Println(tc)
}

func ExampleTestClient_DoJSON() {
	tc := TestClient(func(ctx context.Context, method, relPath string, body any, out any) error {
		// モックの実装例
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, out); err != nil {
			return err
		}
		return nil
	})
	body := map[string]string{"key": "value"}
	out := struct {
		Key string `json:"key"`
	}{
		Key: "value",
	}
	err := tc.DoJSON(context.Background(), "GET", "/example", body, &out)
	fmt.Println(err)
}

package quotetest

import (
	"context"
	"golang/tutorial/todo/internal/quotes"
)

/*
Reference: O'REILLY「実用GO言語」13.4.2 p.299 テストヘルパー用のパッケージを作成する
**_test.goでないため、その他パッケージからimport可能 = 服うすうパッケージにまたがるヘルパー関数を作成可能。
それぞれのパッケージ内の**_test.goからのみimportされればビルド時には含まれないため問題にならない。
あるパッケージの機能をモックする場合は、{package}testという名称が一般的。
*/

// 今回は練習としてquote clientをモックする

// メゾットが一つの場合：

type QuoteClientTest func(ctx context.Context) (quotes.Quote, error)

func (f QuoteClientTest) RandomQuote(ctx context.Context) (quotes.Quote, error) {
	return f(ctx)
}

// メゾットが複数の場合：
// type QuoteClientTest struct {
// 	quotes.HTTPClient
// 	RandomQuoteFunc func(ctx context.Context) (quotes.Quote, error)
// }

// func (c *QuoteClientTest) RandomQuote(ctx context.Context) (quotes.Quote, error) {
// 	if c.RandomQuoteFunc != nil {
// 		return c.RandomQuoteFunc(ctx)
// 	}
// 	return quotes.Quote{}, nil
// }

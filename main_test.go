package main

import (
	"fmt"
	"testing"
)

/*
Reference: O'REILLY「実用GO言語」13.3 p.297 テストに事前ぜ後の処理を追加する
1. ファイル名を"**_test.go"とする
2. テスト関数を"Test**", 引数を"t *testing.T"とする
3. テストの検証
4. テストの実行　```go test -v``` OR ```go test -v -run Test**`
*/
func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}
func setup() {
	fmt.Println("テスト 実行前")
}
func teardown() {
	fmt.Println("テスト 実行後")
}

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	got := Add(1, 2)
	if got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

/*
Reference: O'REILLY「実用GO言語」13.2 p.294 Table Driven Testを実装する
TDT = テストの入力値と期待値をまとめて定義し、テストの実行箇所を一箇所にまとめる
*/
func Calc(a, b int, operator string) (int, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("invalid operator: %s", operator)
	}
}

type Table struct {
	name    string
	args    args
	want    int
	wantErr bool
}
type args struct {
	a  int
	b  int
	op string
}

func TestCalc(t *testing.T) {
	fmt.Println("TestCalc 実行前")

	defer func() {
		fmt.Println("TestCalc 実行後")
	}()

	// テストケース
	tests := []Table{
		{"add ok", args{1, 2, "+"}, 3, false},
		{"sub ok", args{1, 2, "-"}, -1, false},
		{"mul ok", args{1, 2, "*"}, 2, false},
		{"div ok", args{1, 2, "/"}, 0, false},
		{"div ng", args{1, 0, "/"}, 0, true},
		{"invalid op", args{1, 2, "?"}, 0, true},
	}
	for _, tt := range tests {
		fmt.Println("テストケース 実行前")

		defer func() {
			fmt.Println("テストケース 実行後")
		}()
		// 第一引数が各テスト名となる
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calc(tt.args.a, tt.args.b, tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
Reference: O'REILLY「実用GO言語」13.6.4 p.303
時間がかかるテストをスキップしたい場合、```-short```フラグを使う
*/
func TestShortSkip(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

/*
Reference: O'REILLY「実用GO言語」13.11 p.321 ベンチマークを取る
下記の例の実行結果：
```
go test -bench Benchmark -run Benchmark -benchmem

BenchmarkAppendSlice-8            528517              2264 ns/op           25208 B/op         12 allocs/op
BenchmarkAppendSliceBetter-8     1200562               989.6 ns/op          8192 B/op          1 allocs/op
````
*/
func AppendSlice(count, value int) []int {
	var s []int
	for i := 0; i < count; i++ {
		s = append(s, value)
	}
	return s
}
func AppendSliceBetter(count, value int) []int {
	s := make([]int, 0, count)
	for i := 0; i < count; i++ {
		s = append(s, value)
	}
	return s
}
func BenchmarkAppendSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AppendSlice(1000, 42)
	}
}
func BenchmarkAppendSliceBetter(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AppendSliceBetter(1000, 42)
	}
}

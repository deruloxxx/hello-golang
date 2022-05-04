package trace

// byteパッケージはbyte列を持っているだけ。メソッドを好きにつけられる
// https://zenn.dev/hsaki/books/golang-io-package/viewer/bytestring
// ioパッケージは何かに書き込む機能を持つものをまとめて扱うために抽象化されたもの
// https://zenn.dev/hsaki/books/golang-io-package/viewer/io

import (
	"fmt"
	"io"
)

// Tracerはコード内での出来事を記録できるオブジェクトを表すインフェースです
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}
type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// OffはTraceメソッドの呼び出しを無視するTracerを返します。
func Off() Tracer {
	return &nilTracer{}
}

func (t *tracer) Trace(a ...interface{}) {
	// io.writerはユーザーが出力先を自由に選べるからterminalにlogを出力するよう設定
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

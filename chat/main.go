package main

import (
	// text/templateよりhtml/templateの方がセキュア。
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// templは1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// テンプレートファイルのパースを1度だけ実行する
	t.once.Do(func() {
		// mustはerrorが起きた時にnon-nilを返す。ParseFilesでテンプレートをパース。
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// http.Requestにあるポート番号を参照できるようにする
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	// フラグを解釈します(flagの値を取得できるようにする)
	flag.Parse()
	r := newRoom()
	// ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// チャットルームを開始します
	go r.run()
	// Webサーバーを開始します
	log.Println("Webサーバーを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

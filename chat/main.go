package main

import (
	// text/templateよりhtml/templateの方がセキュア。
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
	// nilはtsでいうany型。https://stackoverflow.com/questions/64347531/what-is-the-type-of-nil
	t.templ.Execute(w, nil)
}

func main() {
	// ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})

	// Webサーバーを開始します
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

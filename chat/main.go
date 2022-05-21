package main

import (
	// text/templateよりhtml/templateの方がセキュア。
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// templは1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	gothic.Store = sessions.NewCookieStore([]byte(signature.RandomKey(64)))
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
	)
}

// ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// テンプレートファイルのパースを1度だけ実行する
	t.once.Do(func() {
		// mustはerrorが起きた時にnon-nilを返す。ParseFilesでテンプレートをパース。
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}

	// アプリ用Cookieからユーザー情報を取得する
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	// http.Requestにあるポート番号を参照できるようにする
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":3000", "アプリケーションのアドレス")
	// フラグを解釈します(flagの値を取得できるようにする)
	flag.Parse()
	r := newRoom(UseFileSystemAvatar)
	// ルート
	// authHandlerのServeHttpメソッド実行→認証が成功していたら&templateHandlerのServeHttpメソッド実行
	p := pat.New()
	p.Add("GET", "/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	p.Add("GET", "/login", &templateHandler{filename: "login.html"})
	p.Add("GET", "/upload", &templateHandler{filename: "upload.html"})
	p.Add("GET", "/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	p.Add("GET", "/room", r)
	p.Get("/auth/{provider}/callback", callbackHandler)
	p.Get("/auth/{provider}", gothic.BeginAuthHandler)
	p.Get("/logout", logoutHandler)
	p.Post("/uploader", uploaderHandler)
	// チャットルームを開始します
	go r.run()
	// Webサーバーを開始します
	log.Println("Webサーバーを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, p); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

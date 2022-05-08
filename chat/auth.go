package main

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
)

type authHandler struct {
	// ラップ対象のハンドラを保持
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らかの別のエラーが発生
		panic(err.Error())
	} else {
		// 成功。ラップされたハンドラを呼び出します
		h.next.ServeHTTP(w, r)
	}
}

// ヘルパー関数
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// 外部サービスからの認証結果を判定
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// 外部サービスから取得した情報をアプリ用データとしてCookieにしこむ
	authCookieValue := objx.New(map[string]interface{}{
		"name":       user.UserID,
		"avatar_url": user.AvatarURL,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	// メイン画面へリダイレクト
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: "",
		Path:  "/",
		// MaxAge: -1でブラウザ上のクッキーを即座に削除
		MaxAge: -1,
	})

	w.Header()["Location"] = []string{"/login"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

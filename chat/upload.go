package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	// formのuserIdを取得
	userId := req.FormValue("userId")
	// アップロードされたバイト列を読み込むためのio.Reader型の値を取得。
	// fileにはファイルの値, headerにはファイルに関するメタデータ、errはnilが入る。
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	// バイト列を受け取る
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	// userIdの値を元に保存する画像名を作る。ファイル名の拡張子は元のもの。(headerに含まれている)
	filename := filepath.Join("avatars", userId+filepath.Ext(header.Filename))
	// avatarsフォルダーに新規ファイルを作成してデータを保存。0777は全てのユーザーに対してアクセス券を渡す。
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "成功")
}

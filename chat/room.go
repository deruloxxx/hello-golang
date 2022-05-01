package main

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです。
	// ここで受け取ったメッセージを全て転送する。
	forward chan []byte
}

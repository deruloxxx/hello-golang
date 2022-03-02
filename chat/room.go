package main

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャンネルです。
	forward chan []byte
}

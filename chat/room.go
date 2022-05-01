package main

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです。
	// ここで受け取ったメッセージを全て転送する。
	forward chan []byte
	// 同時アクセスによる競合を防ぐためにjoin, leaveを追加。
	// joinはチャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	// leaveはチャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	// [map]clientsには在室している全てのクライアントをjoin,leaveによって保持。
	clients map[*client]bool
}

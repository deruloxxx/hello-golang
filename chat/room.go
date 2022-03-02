package main

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャンネルです。
	forward chan []byte
	// joinはチャットルームに参加しようとしているクライアントのためのチャンネルです。
	join chan *client
	// leaveはチャットルームから退室しようとしているクライアントのためのチャネルです。
	leave chan *client
	// clientsは在室しているすべてのクライアントが保持されます。
	clients map[*client]bool
}

package main

import (
	"hello-golang/chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです。
	// ここで受け取ったメッセージを全て転送する。
	forward chan *message
	// 同時アクセスによる競合を防ぐためにjoin, leaveを追加。
	// joinはチャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	// leaveはチャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	// [map]clientsには在室している全てのクライアントをjoin,leaveによって保持。
	clients map[*client]bool
	// tracerはチャットルーム上で行われた操作のログを受け取ります。
	tracer trace.Tracer
	// avatarはアバターの情報を取得します。
	avatar Avatar
}

// newRoomはすぐに利用できるチャットルームを生成して返します。
func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		avatar:  avatar,
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました。")
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			// goのチャネルを閉じている
			close(client.send)
			r.tracer.Trace("クライアントが退室しました。")
		case msg := <-r.forward:
			// すべてのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
					r.tracer.Trace("-- クライアントに送信されました。")
				default:
					// 送信に失敗(クリーンアップの処理)
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("-- 送信に失敗しました。クライアントをクリーンアップします。")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// Websocketを使用するためにwebsocket.Upgrader型を使用
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Websocketのコネクションを取得
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました:", err)
		return
	}
	// Websocketのコネクションを取得したらclientを生成してjoinチャネルに渡す
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	// defer文(関数の終了処理)でクライアント終了時に退室。(クリーンアップの処理)
	defer func() { r.leave <- client }()

	// goルーチンの別スレッドとして以下が実行
	go client.write()
	client.read()
}

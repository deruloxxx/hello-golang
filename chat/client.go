package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// clientはチャットを行なっている1人のユーザーを表します。
type client struct {
	// socketはこのクライアントのためのwebsocketです。
	socket *websocket.Conn
	// sendはメッセージが送られるチャネル
	send chan []byte
	// roomはこのクライアントが参加しているチャットルームです。
	room *room
}

//  ReadMessageでデータを読み込む
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			// データをroom.forwardに送信
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// WriteMessageでテータを書き出し
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

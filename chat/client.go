package main

import "github.com/gorilla/websocket"

// *client型はチャットを行なっている1人のユーザーを表します。
type client struct {
	// socketはこのクライアントのためのWebSocketです。(クライアントが通信するためのもの。)
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです。
	send chan []byte
	// roomはこのクライアントが参加しているチャットルームです。
	// Q.room.goを読み込んでいる。room.goを追加するだけでなんで読み込めるかを聞く
	room *room
}

// client型に追加されるメソッド
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			// roomのforwardチャネルに受け取ったmsgを送信
			// <-が送信
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	// for ...rangeはforeach的なやつ
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

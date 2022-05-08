package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// *client型はチャットを行なっている1人のユーザーを表します。
type client struct {
	// socketはこのクライアントのためのWebSocketです。(クライアントが通信するためのもの。)
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです。
	send chan *message
	// roomはこのクライアントが参加しているチャットルームです。
	// Q.room.goを読み込んでいる。room.goを追加するだけでなんで読み込めるかを聞く
	room *room
	// userDataはユーザーに関する情報を保持します
	userData map[string]interface{}
}

// client型に追加されるメソッド
func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
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
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

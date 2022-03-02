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

// newRoomはすぐに利用できるチャットルームを生成して返します。
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// すべてのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

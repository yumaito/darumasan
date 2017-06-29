package app

import (
	"log"
)

// Hub は登録されているクライアントの管理やメッセージのやり取りを管理する中枢の役割を果たします
type Hub struct {
	clients    map[*Client]bool
	message    chan *Message
	register   chan *Client
	unregister chan *Client
	write      chan []byte
	logger     *log.Logger
}

func NewHub(logger *log.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		message:    make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		write:      make(chan []byte),
		logger:     logger,
	}
}

// run は各chanからの入力に対して処理を行います
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.logger.Printf("%+v connected\n", client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.write)
				h.logger.Printf("%+v disconnected\n", client)
			}
		case msg := <-h.message:
			h.logger.Printf("%+v\n", msg)
			switch msg.Type {
			case CLIENT_TYPE_CLIENT:
				h.send(CLIENT_TYPE_CURATOR, msg.Msg)
			case CLIENT_TYPE_CURATOR:
				h.send(CLIENT_TYPE_CLIENT, msg.Msg)
			}
			// case clientEvent := <-h.clientEvent:
			// スマホクライアントからの入力
			// case curatorEvent := <-h.curatorEvent:
			// 鬼クライアントからの入力
		}
	}
}

// send は指定したclientTypeのclient全てにメッセージを送ります
func (h *Hub) send(clientType uint32, msg []byte) {
	for client, _ := range h.clients {
		if client.clientType == clientType {
			select {
			case client.write <- msg:
			}
		}
	}
}

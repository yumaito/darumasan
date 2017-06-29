package app

import (
	"log"
)

// Hub は登録されているクライアントの管理やメッセージのやり取りを管理する中枢の役割を果たします
type Hub struct {
	clients      map[*Client]bool
	clientEvent  chan *ClientEvent
	curatorEvent chan *CuratorEvent
	register     chan *Client
	unregister   chan *Client
	logger       *log.Logger
}

func NewHub(logger *log.Logger) *Hub {
	return &Hub{
		clients:      make(map[*Client]bool),
		clientEvent:  make(chan *ClientEvent),
		curatorEvent: make(chan *CuratorEvent),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		logger:       logger,
	}
}

// run は各chanからの入力に対して処理を行います
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case clientEvent := <-h.clientEvent:
			// スマホクライアントからの入力
		case curatorEvent := <-h.curatorEvent:
			// 鬼クライアントからの入力
		}
	}
}

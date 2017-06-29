package app

import (
	"log"
)

// Hub は登録されているクライアントの管理やメッセージのやり取りを管理する中枢の役割を果たします
type Hub struct {
	clients    map[*Client]bool
	curatorID  string
	message    chan *ClientMessage
	broadcast  chan *GameMessage
	register   chan *Client
	unregister chan *Client
	IsWatched  bool
	logger     *log.Logger
}

func NewHub(logger *log.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		message:    make(chan *ClientMessage),
		broadcast:  make(chan *GameMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

// run は各chanからの入力に対して処理を行います
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.logger.Printf("id:%s type:%d connected\n", client.ID, client.clientType)
			if client.clientType == CLIENT_TYPE_CURATOR {
				h.curatorID = client.ID
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.write)
				h.logger.Printf("id:%s type:%d disconnected\n", client.ID, client.clientType)
			}
		case msg := <-h.message:
			h.logger.Printf("receive from %s %+v\n", msg.ID, msg)
			gm := h.complementMessage(msg)
			h.send(msg.ID, gm)
		}
	}
}

func (h *Hub) complementMessage(cm *ClientMessage) *GameMessage {
	cs := make([]string, 0)
	for key, _ := range h.clients {
		cs = append(cs, key.ID)
	}
	switch cm.ClientType {
	case CLIENT_TYPE_CLIENT:
		// クライアントからのメッセージ
	case CLIENT_TYPE_CURATOR:
		// 鬼からのメッセージならis_watchedを更新
		h.IsWatched = cm.Status
	}

	return &GameMessage{
		Clients:     cs,
		DeadClients: []string{},
		CuratorID:   h.curatorID,
		IsWatched:   h.IsWatched,
	}
}

func (h *Hub) send(cid string, gm *GameMessage) {
	// 送り主と鬼にメッセージを送信
	for key, _ := range h.clients {
		if key.ID == cid || key.ID == h.curatorID {
			h.logger.Printf("send %+v\n", gm)
			key.write <- gm
		}
	}
}

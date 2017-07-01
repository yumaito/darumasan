package app

import (
	"time"

	"go.uber.org/zap"
)

// Hub は登録されているクライアントの管理やメッセージのやり取りを管理する中枢の役割を果たします
type Hub struct {
	clients     map[*Client]bool
	curatorID   string
	deadClients []string
	message     chan *ClientMessage
	register    chan *Client
	unregister  chan *Client
	IsWatched   bool
	logger      *zap.Logger
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		deadClients: []string{},
		message:     make(chan *ClientMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		logger:      logger,
	}
}

// run は各chanからの入力に対して処理を行います
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// 接続クライアントの登録
			// 鬼が接続した場合は鬼として登録
			h.clients[client] = true
			h.logger.Info("client connected",
				zap.String("client_id", client.ID),
				zap.Uint32("client_type", client.clientType),
			)
			if client.clientType == CLIENT_TYPE_CURATOR {
				h.curatorID = client.ID
			}
			initMsg := h.messageByConnectionEvent(client)
			h.sendToClientAndCurator(client.ID, initMsg)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.write)
				// クライアントが切断した場合鬼側にも送る
				if client.clientType == CLIENT_TYPE_CLIENT {
					msg := h.messageByConnectionEvent(client)
					h.sendToClientAndCurator(h.curatorID, msg)
				}
				h.logger.Info("client disconnected",
					zap.String("client_id", client.ID),
					zap.Uint32("client_type", client.clientType),
				)
			}
		case msg := <-h.message:
			gm := h.createMessage(msg)
			h.sendToClientAndCurator(msg.ID, gm)
		}
	}
}

func (h *Hub) createMessage(cm *ClientMessage) *GameMessage {
	cs := make([]string, 0)
	for key, _ := range h.clients {
		cs = append(cs, key.ID)
	}
	switch cm.ClientType {
	case CLIENT_TYPE_CLIENT:
		// クライアントからのメッセージ
		// 監視中にtrueがきたらアウトにする
		// ユニークにするため一度mapに入れて整理している
		if h.IsWatched && cm.Status {
			h.deadClients = append(h.deadClients, cm.ID)
			m := make(map[string]bool)
			for _, val := range h.deadClients {
				m[val] = true
			}
			m[cm.ID] = true
			dc := make([]string, 0)
			for key, _ := range m {
				dc = append(dc, key)
			}
			h.deadClients = dc
		}
	case CLIENT_TYPE_CURATOR:
		// 鬼からのメッセージならis_watchedを更新
		h.IsWatched = cm.Status
	}

	return &GameMessage{
		From:        cm.ID,
		ClientType:  cm.ClientType,
		Clients:     cs,
		DeadClients: h.deadClients,
		CuratorID:   h.curatorID,
		IsWatched:   h.IsWatched,
		CreatedAt:   time.Now(),
	}
}

// messageByConnectionEvent はクライアントの接続に応じてメッセージを作成します
func (h *Hub) messageByConnectionEvent(client *Client) *GameMessage {
	cs := make([]string, 0)
	for key, _ := range h.clients {
		cs = append(cs, key.ID)
	}
	return &GameMessage{
		From:        client.ID,
		ClientType:  client.clientType,
		Clients:     cs,
		DeadClients: h.deadClients,
		CuratorID:   h.curatorID,
		IsWatched:   h.IsWatched,
		CreatedAt:   time.Now(),
	}
}

// sendToClientAndCurator はcidのクライアントと鬼に対してメッセージを送信します
func (h *Hub) sendToClientAndCurator(cid string, gm *GameMessage) {
	// 送り主と鬼にメッセージを送信
	for key, _ := range h.clients {
		if key.ID == cid || key.ID == h.curatorID {
			msg := &GameMessage{
				From:        gm.From,
				To:          key.ID,
				ClientType:  gm.ClientType,
				Clients:     gm.Clients,
				DeadClients: gm.DeadClients,
				CuratorID:   gm.CuratorID,
				IsWatched:   gm.IsWatched,
				CreatedAt:   gm.CreatedAt,
			}
			key.write <- msg
		}
	}
}

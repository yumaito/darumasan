package app

import (
	"time"
)

// GameMessage はサーバーから送信されるメッセージ
type GameMessage struct {
	From        string    `json:"from"`
	To          string    `json:"to"`
	ClientType  uint32    `json:"client_type"`
	Clients     []string  `json:"clients"`
	DeadClients []string  `json:"dead_clients"`
	CuratorID   string    `json:"curator_id"`
	IsWatched   bool      `json:"is_watched"`
	CreatedAt   time.Time `json:"created_at"`
}

// Message はクライアントから直接送られてくるメッセージ
type Message struct {
	Status bool `json:"status"`
}

// ClientMessage はクライアントから送られてきたメッセージにクライアントの情報を負荷したもの
// Hubへの送信に使われる
type ClientMessage struct {
	ID         string `json:"id"`
	ClientType uint32 `json:"client_type"`
	Status     bool   `json:"status"`
}

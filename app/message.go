package app

import (
	"time"

	"go.uber.org/zap/zapcore"
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

// MarshallLogObject はzapによるログ用のメソッドです
func (g *GameMessage) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("from", g.From)
	enc.AddString("to", g.To)
	enc.AddUint32("client_type", g.ClientType)
	enc.AddArray("clients", Strings(g.Clients))
	enc.AddArray("dead_clients", Strings(g.DeadClients))
	enc.AddString("curator_id", g.CuratorID)
	enc.AddBool("is_watched", g.IsWatched)
	enc.AddTime("created_at", g.CreatedAt)

	return nil
}

type Strings []string

func (ss Strings) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, s := range ss {
		enc.AppendString(s)
	}
	return nil
}

// Message はクライアントから直接送られてくるメッセージ
type Message struct {
	Status bool `json:"status"`
}

func (m *Message) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("status", m.Status)
	return nil
}

// ClientMessage はクライアントから送られてきたメッセージにクライアントの情報を負荷したもの
// Hubへの送信に使われる
type ClientMessage struct {
	ID         string `json:"id"`
	ClientType uint32 `json:"client_type"`
	Status     bool   `json:"status"`
}

func (c *ClientMessage) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", c.ID)
	enc.AddUint32("client_type", c.ClientType)
	enc.AddBool("status", c.Status)
	return nil
}

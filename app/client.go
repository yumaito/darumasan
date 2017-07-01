package app

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	CLIENT_TYPE_CLIENT  = 1
	CLIENT_TYPE_CURATOR = 2
	CLIENT_TYPE_BUTTON  = 3
	CLIENT_TYPE_MONITOR = 4
)

type Client struct {
	hub        *Hub
	ID         string
	conn       *websocket.Conn
	write      chan *GameMessage
	clientType uint32
}

func NewClient(hub *Hub, conn *websocket.Conn, id string, clientType uint32) *Client {
	return &Client{
		hub:        hub,
		ID:         id,
		conn:       conn,
		write:      make(chan *GameMessage),
		clientType: clientType,
	}
}

func (c *Client) Disconnect() {
	c.conn.Close()
	c.hub.unregister <- c
}

func (c *Client) ReadMessage() {
	defer c.Disconnect()
	for {
		m := &Message{}
		if err := c.conn.ReadJSON(m); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				c.hub.logger.Error("ReadJSON",
					zap.String("msg", err.Error()),
				)
			}
			break
		}
		cm := &ClientMessage{
			ID:         c.ID,
			ClientType: c.clientType,
			Status:     m.Status,
		}
		c.hub.logger.Info("message received",
			zap.Object("message", m),
		)
		c.hub.message <- cm
	}
}

func (c *Client) WriteMessage() {
	defer c.Disconnect()
	for {
		select {
		case message, ok := <-c.write:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				c.hub.logger.Error("WriteJSON",
					zap.String("msg", err.Error()),
				)
				return
			}
			c.hub.logger.Info("send message",
				zap.Object("message", message),
			)
		}
	}
}

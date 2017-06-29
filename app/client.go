package app

import (
	"github.com/gorilla/websocket"
)

const (
	CLIENT_TYPE_CLIENT  = 1
	CLIENT_TYPE_CURATOR = 2
)

var (
	newline = []byte{'\n'}
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
	c.hub.logger.Println("Disconnect:", c.ID)
}

func (c *Client) ReadMessage() {
	defer c.Disconnect()
	for {
		cm := &ClientMessage{}
		if err := c.conn.ReadJSON(cm); err != nil {
			c.hub.logger.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.hub.logger.Println(err)
			}
			break
		}
		c.hub.message <- cm
	}
}

func (c *Client) WriteMessage() {
	defer c.Disconnect()
	for {
		select {
		case message, ok := <-c.write:
			c.hub.logger.Printf("client_id:%s %+v\n", c.ID, message)
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				c.hub.logger.Println(err)
			}
		}
	}
}

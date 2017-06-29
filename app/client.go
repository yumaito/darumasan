package app

import (
	"bytes"
	"io"

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
	conn       *websocket.Conn
	write      chan []byte
	clientType uint32
}

func NewClient(hub *Hub, conn *websocket.Conn, clientType uint32) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		write:      make(chan []byte),
		clientType: clientType,
	}
}

func (c *Client) Disconnect() {
	c.conn.Close()
	c.hub.unregister <- c
	c.hub.logger.Println("Disconnect")
}

func (c *Client) ReadMessage() {
	defer c.Disconnect()
	for {
		_, reader, err := c.conn.NextReader()
		if err != nil {
			c.hub.logger.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.hub.logger.Println(err)
			}
			break
		}
		var b []byte
		w := bytes.NewBuffer(b)
		if _, err := io.Copy(w, reader); err != nil {
			c.hub.logger.Println(err)
		}
		m := NewMessage(c.clientType, w.Bytes())
		c.hub.message <- m
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
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.hub.logger.Println(err)
			}
			w.Write(message)
			// キューされたmessageを順に処理
			n := len(c.write)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.write)
			}
			if err := w.Close(); err != nil {
				c.hub.logger.Println(err)
			}
		}
	}
}

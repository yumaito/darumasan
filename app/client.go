package app

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	CLIENT_TYPE_CLIENT  = 1
	CLIENT_TYPE_CURATOR = 2
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan []*ClientEvent
	clientType uint32
}

func NewClient(hub *Hub, conn *websocket.Conn, clientType uint32) {
	return &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan send),
		clientType: clientType,
	}
}

func WsServer(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.logger.Println(err)
		return
	}
	client := NewClient(hub, conn)
}

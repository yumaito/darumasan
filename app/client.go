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

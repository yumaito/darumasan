package app

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ClientHandler はスマホ端末側がつなぐエンドポイントのhandlerです
func ClientHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.logger.Println("error upgrade:", err)
		return
	}
	client := NewClient(hub, conn, CLIENT_TYPE_CLIENT)

	hub.register <- client
	hub.logger.Println("connected /client")
	go client.ReadMessage()
	go client.WriteMessage()
}

// CuratorHandler は鬼側がつなぐエンドポイントのhandlerです
func CuratorHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.logger.Println("error upgrade:", err)
	}
	client := NewClient(hub, conn, CLIENT_TYPE_CURATOR)
	hub.register <- client
	hub.logger.Println("connected /curator")
	go client.ReadMessage()
	go client.WriteMessage()
}

package app

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func generateID() string {
	h := sha256.New()
	str := strconv.FormatInt(time.Now().UnixNano(), 10)
	io.WriteString(h, str)
	result := fmt.Sprintf("%x", h.Sum(nil))
	return result[:10]
}

// ClientHandler はスマホ端末側がつなぐエンドポイントのhandlerです
func ClientHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.logger.Error("client connection",
			zap.String("msg", err.Error()),
		)
		return
	}
	id := generateID()
	client := NewClient(hub, conn, id, CLIENT_TYPE_CLIENT)

	hub.register <- client
	go client.ReadMessage()
	go client.WriteMessage()
}

// CuratorHandler は鬼側がつなぐエンドポイントのhandlerです
func CuratorHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.logger.Error("curator connection",
			zap.String("msg", err.Error()),
		)
	}
	id := generateID()
	client := NewClient(hub, conn, id, CLIENT_TYPE_CURATOR)
	hub.register <- client
	go client.ReadMessage()
	go client.WriteMessage()
}

func ButtonHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.logger.Error("button connection",
			zap.String("msg", err.Error()),
		)
	}
	id := generateID()
	client := NewClient(hub, conn, id, CLIENT_TYPE_BUTTON)
	hub.register <- client
	go client.ReadMessage()
	go client.WriteMessage()
}

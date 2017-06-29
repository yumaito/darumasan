package app

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Handler is XXX
type Handler struct {
	logger      *log.Logger
	upgrader    websocket.Upgrader
	clientChan  chan ClientEvent
	curatorChan chan CuratorEvent
}

// NewHandler is XXX
func NewHandler(logger *log.Logger) *Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Handler{
		logger:      logger,
		upgrader:    upgrader,
		clientChan:  make(chan ClientEvent),
		curatorChan: make(chan CuratorEvent),
	}
}

// ClientHandler はスマホ端末側がつなぐエンドポイントのhandlerです
func (h *Handler) ClientHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Println("error upgrade:", err)
		return
	}
	client := NewClient(hub, conn, CLIENT_TYPE_CLIENT)
	h.logger.Println("connected /client")
	defer conn.Close()
	defer h.logger.Println("disconnected /client")

}

// CuratorHandler は鬼側がつなぐエンドポイントのhandlerです
func (h *Handler) CuratorHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Println("error upgrade:", err)
	}
	h.logger.Println("connected /curator")
	defer conn.Close()
	defer h.logger.Println("disconnected /curator")

}

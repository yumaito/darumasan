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

// EchoHandler は入力されたものをそのまま返すhandlerです
func (h *Handler) EchoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Println("error upgrade:", err)
		return
	}
	h.logger.Println("connected /echo")
	defer conn.Close()
	defer h.logger.Println("disconnected /echo")
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			h.logger.Println("error next reader:", err)
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			h.logger.Println("error next writer:", err)
			return
		}

		if _, err := io.Copy(w, r); err != nil {
			h.logger.Println("error copy:", err)
			return
		}
		if err := w.Close(); err != nil {
			h.logger.Println("error close:", err)
			return
		}
	}
}

// ClientHandler はスマホ端末側がつなぐエンドポイントのhandlerです
func (h *Handler) ClientHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Println("error upgrade:", err)
	}
	h.logger.Println("connected /client")
	defer conn.Close()
	defer h.logger.Println("disconnected /client")

	for {
		e := NewClientEvent()
		if err := conn.ReadJSON(e); err != nil {
			h.logger.Println("error readJSON:", err)
			return
		}
		h.logger.Println("/client received:", r)
		h.clientChan <- *e

		select {
		case e := <-h.curatorChan:
			if err := conn.WriteJSON(&e); err != nil {
				h.logger.Println("error writeJSON:", err)
				return
			}
		}
	}
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

	for {
		e := NewCuratorEvent()
		if err := conn.ReadJSON(e); err != nil {
			h.logger.Println("error readJSON", err)
			return
		}
		h.logger.Println("/curator received:", r)
		h.curatorChan <- *e

		select {
		case e := <-h.clientChan:
			if err := conn.WriteJSON(&e); err != nil {
				h.logger.Println("error writeJSON:", err)
				return
			}
		}
	}
}

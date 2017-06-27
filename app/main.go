package app

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrade:", err)
		return
	}
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Println("error next reader:", err)
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			log.Println("error next writer:", err)
			return
		}

		if _, err := io.Copy(w, r); err != nil {
			log.Println("error copy:", err)
			return
		}
		if err := w.Close(); err != nil {
			log.Println("error close:", err)
			return
		}
	}
}

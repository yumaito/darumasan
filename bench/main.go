package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yumaito/darumasan/app"
)

const (
	SERVER            = "localhost:8080"
	CLIENT_END_POINT  = "client"
	CURATOR_END_POINT = "curator"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	logger := log.New(os.Stdout, "", log.Lshortfile)

	u := url.URL{
		Scheme: "ws",
		Host:   SERVER,
		Path:   "/" + CLIENT_END_POINT,
	}
	logger.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Fatal("dial:", err)
	}
	defer c.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m := &app.Message{
				Status: true,
			}
			if err := c.WriteJSON(m); err != nil {
				logger.Println(err)
				return
			}
			logger.Printf("send message: %+v\n", m)
		case <-interrupt:
			logger.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Println("write close:", err)
				return
			}
			c.Close()
			return
		}
	}
}

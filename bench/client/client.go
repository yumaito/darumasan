package client

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yumaito/darumasan/app"
)

const (
	SERVER            = "localhost:8080"
	CLIENT_END_POINT  = "client"
	CURATOR_END_POINT = "curator"
	BUTTON_END_POINT  = "button"
)

type Client struct {
	rate         uint32
	tickDuration time.Duration
	clientType   uint32
	interrupt    chan os.Signal
	logger       *log.Logger
	url          url.URL
}

func NewClient(rate uint32, duration time.Duration, clientType uint32, logger *log.Logger) (*Client, error) {
	var path string
	switch clientType {
	case app.CLIENT_TYPE_CLIENT:
		path = "/" + CLIENT_END_POINT
	case app.CLIENT_TYPE_CURATOR:
		path = "/" + CURATOR_END_POINT
	case app.CLIENT_TYPE_BUTTON:
		path = "/" + BUTTON_END_POINT
	default:
		return nil, fmt.Errorf("invalid clientType:%d", clientType)
	}

	u := url.URL{
		Scheme: "ws",
		Host:   SERVER,
		Path:   path,
	}
	return &Client{
		rate:         rate,
		tickDuration: duration,
		clientType:   clientType,
		logger:       logger,
		url:          u,
	}, nil
}

func (c *Client) Run(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(c.url.String(), nil)
	if err != nil {
		c.logger.Println("dial:", err)
		return
	}
	c.logger.Printf("connected to:%s\n", c.url.String())
	defer conn.Close()
	defer c.logger.Printf("disconnected from:%s\n", c.url.String())
	for {
		select {
		case <-ctx.Done():
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				c.logger.Println("write close:", err)
				return
			}
			conn.Close()
			return
		}
	}
}

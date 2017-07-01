package client

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yumaito/darumasan/app"
)

const (
	CLIENT_END_POINT  = "client"
	CURATOR_END_POINT = "curator"
	BUTTON_END_POINT  = "button"
)

type Client struct {
	config *Config
	logger *log.Logger
	url    url.URL
}

type Config struct {
	Rate     int    `yaml:"rate"`
	Duration int64  `yaml:"duration"`
	Type     uint32 `yaml:"type"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewClient(conf *Config, host string, logger *log.Logger) (*Client, error) {
	var path string
	switch conf.Type {
	case app.CLIENT_TYPE_CLIENT:
		path = "/" + CLIENT_END_POINT
	case app.CLIENT_TYPE_CURATOR:
		path = "/" + CURATOR_END_POINT
	case app.CLIENT_TYPE_BUTTON:
		path = "/" + BUTTON_END_POINT
	default:
		return nil, fmt.Errorf("invalid clientType:%d", conf.Type)
	}

	u := url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   path,
	}
	return &Client{
		config: conf,
		logger: logger,
		url:    u,
	}, nil
}

func (c *Client) Run(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(c.url.String(), nil)
	if err != nil {
		c.logger.Println("dial:", err)
		return
	}
	c.logger.Printf("connected to:%s\n", c.url.String())
	t := time.Duration(c.config.Duration) * time.Millisecond
	ticker := time.NewTicker(t)

	status := false
	defer conn.Close()
	defer c.logger.Printf("disconnected from:%s\n", c.url.String())
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 鬼は一定時間ごとにtrueとfalseをトグルするだけ
			switch c.config.Type {
			case app.CLIENT_TYPE_CURATOR:
				status = true
			default:
				status = c.rateSelector(c.config.Rate)
			}
			if c.config.Type == app.CLIENT_TYPE_CURATOR {
				status = !status
			} else {
				status = true
			}
			m := &app.Message{
				Status: status,
			}
			if err := conn.WriteJSON(m); err != nil {
				c.logger.Println(err)
				return
			}
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

func (c *Client) rateSelector(rate int) bool {
	number := rand.Intn(100)
	if (number + rate) >= 100 {
		return true
	}
	return false
}

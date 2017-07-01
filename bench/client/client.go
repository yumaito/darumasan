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
	config      *Config
	logger      *log.Logger
	url         url.URL
	sendStopped bool
	status      bool
	conn        *websocket.Conn
}

type Config struct {
	Rate     int    `yaml:"rate"`     // trueになる確率(playerの場合アウトになる確率)
	Duration int64  `yaml:"duration"` // 信号を送信するインターバル(ミリ秒)
	Type     uint32 `yaml:"type"`     // クライントのタイプ
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
	c := &Client{
		config:      conf,
		logger:      logger,
		url:         u,
		sendStopped: false,
		status:      false,
	}
	return c, nil
}

func (c *Client) Stop() {
	c.conn.Close()
	c.logger.Printf("disconnected from:%s\n", c.url.String())
}

func (c *Client) Run(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(c.url.String(), nil)
	if err != nil {
		c.logger.Println("dial:", err)
		return
	}
	c.conn = conn
	c.logger.Printf("connected to:%s\n", c.url.String())
	t := time.Duration(c.config.Duration) * time.Millisecond
	ticker := time.NewTicker(t)

	go c.read(ctx)
	defer c.Stop()
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 鬼は一定時間ごとにtrueとfalseをトグルするだけ
			// 送信ストップフラグが立っていないときだけ送信する
			//（アウトになったときにフラグが立つ）
			if !c.sendStopped {
				switch c.config.Type {
				case app.CLIENT_TYPE_CURATOR:
					c.status = !c.status
					m := &app.Message{
						Status: c.status,
					}
					if err := conn.WriteJSON(m); err != nil {
						c.logger.Println(err)
						return
					}
					c.logger.Printf("send to %s, status: %v", c.url.String(), m.Status)
				default:
					if c.rateSelector(c.config.Rate) {
						m := &app.Message{
							Status: true,
						}
						if err := conn.WriteJSON(m); err != nil {
							c.logger.Println(err)
							return
						}
						c.logger.Printf("send to %s, status: %v", c.url.String(), m.Status)
					}
				}
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

func (c *Client) read(ctx context.Context) {
	defer c.conn.Close()
	for {
		m := &app.GameMessage{}
		if err := c.conn.ReadJSON(m); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				c.logger.Println(err)
			}
			return
		}
		if !c.sendStopped {
			for _, id := range m.DeadClients {
				if id == m.To {
					c.sendStopped = true
				}
			}
			if c.sendStopped {
				c.logger.Printf("connection to %s: send stopped", c.url.String())
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
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

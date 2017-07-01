package main

import (
	"context"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/Code-Hex/sigctx"
	"github.com/yumaito/darumasan/app"
	"github.com/yumaito/darumasan/bench/client"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)

	clients := make([]*client.Client, 5)
	for i := 0; i < len(clients); i++ {
		c, err := client.NewClient(100, time.Second, app.CLIENT_TYPE_CLIENT, logger)
		if err != nil {
			logger.Println(err)
			return
		}
		clients[i] = c
	}
	curator, err := client.NewClient(100, 3*time.Second, app.CLIENT_TYPE_CURATOR, logger)
	if err != nil {
		logger.Println(err)
		return
	}
	clients = append(clients, curator)

	wg := &sync.WaitGroup{}
	ctx := sigctx.WithCancelSignals(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	for _, c := range clients {
		wg.Add(1)
		go func(c *client.Client) {
			logger.Printf("%p: %+v\n", c, c)
			c.Run(ctx)
			wg.Done()
		}(c)
	}
	wg.Wait()
}

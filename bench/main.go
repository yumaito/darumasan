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

	wg := &sync.WaitGroup{}
	ctx := sigctx.WithCancelSignals(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	for _, client := range clients {
		wg.Add(1)
		go func() {
			client.Run(ctx)
			wg.Done()
		}()
	}
	wg.Wait()

	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()
	//for {
	//	select {
	//	case <-ticker.C:
	//		m := &app.Message{
	//			Status: true,
	//		}
	//		if err := c.WriteJSON(m); err != nil {
	//			logger.Println(err)
	//			return
	//		}
	//		logger.Printf("send message: %+v\n", m)
	//	case <-interrupt:
	//		logger.Println("interrupt")
	//		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	//		if err != nil {
	//			logger.Println("write close:", err)
	//			return
	//		}
	//		c.Close()
	//		return
	//	}
	//}
}

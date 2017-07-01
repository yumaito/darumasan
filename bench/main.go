package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"syscall"

	"github.com/Code-Hex/sigctx"
	"github.com/yumaito/darumasan/bench/client"
	"gopkg.in/yaml.v2"
)

var (
	configFile = flag.String("c", "config.yml", "config file")
)

type Config struct {
	Host    string `yaml:"host"`
	Clients []struct {
		Number int            `yaml:"number"`
		Client *client.Config `yaml:"client"`
	} `yaml:"clients"`
}

func main() {
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Lshortfile)
	// configの読み込み
	conf := &Config{}
	buf, err := ioutil.ReadFile(*configFile)
	if err != nil {
		logger.Println(err)
		return
	}
	if err := yaml.Unmarshal(buf, conf); err != nil {
		logger.Println(err)
		return
	}
	//
	clients := make([]*client.Client, 0)
	for _, c := range conf.Clients {
		for i := 0; i < c.Number; i++ {
			client, err := client.NewClient(c.Client, conf.Host, logger)
			if err != nil {
				logger.Println(err)
				return
			}
			clients = append(clients, client)
		}
	}

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
			c.Run(ctx)
			wg.Done()
		}(c)
	}
	wg.Wait()
}

package main

import (
	"log"
	"net/http"

	"github.com/yumaito/darumasan/app"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println("error new logger:", err)
		return
	}
	logger.Info("server started")

	hub := app.NewHub(logger)
	go hub.Run()
	http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		app.ClientHandler(hub, w, r)
	})
	http.HandleFunc("/curator", func(w http.ResponseWriter, r *http.Request) {
		app.CuratorHandler(hub, w, r)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("ListenAndServe",
			zap.String("msg", err.Error()),
		)
	}
}

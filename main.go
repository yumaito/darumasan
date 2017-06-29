package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yumaito/darumasan/app"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	handler := app.NewHandler(logger)

	hub := app.NewHub(logger)
	go hub.Run()
	http.HandleFunc("/client", handler.ClientHandler)
	http.HandleFunc("/curator", handler.CuratorHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Println("ListenAndServe:", err.Error())
	}
}

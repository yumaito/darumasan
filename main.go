package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yumaito/darumasan/app"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	logger.Println("server started now")

	hub := app.NewHub(logger)
	go hub.Run()
	http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		app.ClientHandler(hub, w, r)
	})
	http.HandleFunc("/curator", func(w http.ResponseWriter, r *http.Request) {
		app.CuratorHandler(hub, w, r)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Println("ListenAndServe:", err.Error())
	}
}

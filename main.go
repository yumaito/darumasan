package main

import (
	"log"
	"net/http"

	"github.com/yumaito/darumasan/app"
)

func main() {
	http.HandleFunc("/echo", app.EchoHandler)
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Println("ListenAndServer:", err.Error())
	}
}

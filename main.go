package main

import (
	handlers "calc/backend"
	"fmt"
	"net/http"
)

func main() {

	handlers.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/logExp", handlers.LogHandler)
	mux.HandleFunc("/chat", handlers.ChatHandler)
	mux.HandleFunc("/sendMsg", handlers.ChatPostHandler)
	mux.HandleFunc("/getMsgs", handlers.GetMsgsHandler)
	mux.HandleFunc("/ws", handlers.WsHandler)

	mux.HandleFunc("GET /static/", handlers.StaticHandler)

	fmt.Println("Started http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

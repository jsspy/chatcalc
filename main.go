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
	mux.HandleFunc("POST /logExp", handlers.LogHandler)
	mux.HandleFunc("GET /chat", handlers.ChatHandler)
	mux.HandleFunc("POST /sendMsg", handlers.ChatPostHandler)
	mux.HandleFunc("GET /getMsgs", handlers.GetMsgsHandler)

	mux.HandleFunc("GET /static/", handlers.StaticHandler)

	fmt.Println("Started http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

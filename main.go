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
	mux.HandleFunc("/chat", handlers.ChatHandler)
	mux.HandleFunc("/ws", handlers.WsHandler)
	mux.HandleFunc("/api/upload-image", handlers.UploadImageHandler)

	// Serve frontend static files
	mux.HandleFunc("/static/", handlers.StaticHandler)

	// Serve uploaded files from ./uploads at /uploads/
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	fmt.Println("Started http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

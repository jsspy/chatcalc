package main

import (
	handlers "calc/backend"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize database
	if err := handlers.InitDB("./chat.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}
	defer handlers.CloseDB()

	// Static file server for uploads
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Endpoints
	http.Handle("/", http.FileServer(http.Dir("./frontend/")))
	http.HandleFunc("/ws", handlers.HandleWS)
	http.HandleFunc("/upload", handlers.HandleUpload)
	http.HandleFunc("/messages", handlers.HandleGetMessages)

	fmt.Println("Listening on http://localhost:8080/chat.html")

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down...")
		handlers.CloseDB()
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

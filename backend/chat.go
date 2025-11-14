package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./frontend/chat.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	messages, err := GetAllChatMessages()
	if err != nil {
		log.Printf("error getting chat messages: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, messages); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

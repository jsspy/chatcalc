package handlers

import (
	"encoding/json"
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

func GetMsgsHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := GetAllChatMessages()
	if err != nil {
		log.Printf("error getting chat messages: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": "true",
		"posts":   messages,
	})
}

func ChatPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var postData struct {
		Author  string `json:"user"`
		Message string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Invalid JSON",
		})
		return
	}

	if postData.Message == "" {
		http.Error(w, "message cannot be empty", http.StatusBadRequest)
		return
	}

	if err := SaveChatMessage(postData.Message, postData.Author); err != nil {
		log.Printf("error saving message: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

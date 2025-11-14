package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := GetMessages()
	if err != nil {
		log.Println("Error retrieving messages:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve messages"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if messages == nil {
		messages = []Message{}
	}
	json.NewEncoder(w).Encode(messages)
}

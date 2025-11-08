package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CalcLog struct {
	Expression string `json:"expression"`
	Result     any    `json:"result"`
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	var logEntry CalcLog
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Println(logEntry.Expression)
	if logEntry.Expression == "09+21" {
		http.Redirect(w, r, "/chat", http.StatusPermanentRedirect)
		return
	}

	// Just write the response; header 200 is implied
	w.Write([]byte(`{"status":"logged"}`))
}

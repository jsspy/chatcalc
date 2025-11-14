package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	Text      string `json:"text"`
	FileURL   string `json:"fileUrl"`
	ReplyToID string `json:"replyToId"`
}

var clients = make(map[*websocket.Conn]string) // conn → username (A or J)
var clientsMu sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Detect device → user
func detectUser(r *http.Request) string {
	ua := strings.ToLower(r.UserAgent())
	if strings.Contains(ua, "iphone") {
		return "J"
	}
	return "A"
}

// Generate unique ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func HandleWS(w http.ResponseWriter, r *http.Request) { // Broadcast to all clients
	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	user := detectUser(r)

	clientsMu.Lock()
	clients[conn] = user
	clientsMu.Unlock()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse message and save to database
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		// Generate unique ID if not provided
		if msg.ID == "" {
			msg.ID = generateID()
		}

		// Save to SQLite
		if err := SaveMessage(&msg); err != nil {
			log.Println("Error saving message to database:", err)
			continue
		}

		// Re-marshal with the ID for broadcasting
		data, _ = json.Marshal(msg)

		// Forward message to both users
		broadcast(data)
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	conn.Close()
}
func broadcast(data []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for c := range clients {
		c.WriteMessage(websocket.TextMessage, data)
	}
}

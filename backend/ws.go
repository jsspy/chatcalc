package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Message envelope used over websocket
type wsEnvelope struct {
	Type  string        `json:"type"`
	Post  *ChatMessage  `json:"post,omitempty"`
	Posts []ChatMessage `json:"posts,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub maintains active connections and broadcasts messages.
type Hub struct {
	// Registered connections.
	conns map[*websocket.Conn]bool
	// Inbound messages from the connections.
	broadcast chan wsEnvelope
	// Register requests from the connections.
	register chan *websocket.Conn
	// Unregister requests from connections.
	unregister chan *websocket.Conn
}

var hub = Hub{
	conns:      make(map[*websocket.Conn]bool),
	broadcast:  make(chan wsEnvelope),
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.conns[c] = true
		case c := <-h.unregister:
			if _, ok := h.conns[c]; ok {
				delete(h.conns, c)
				c.Close()
			}
		case msg := <-h.broadcast:
			// broadcast to all conns
			b, err := json.Marshal(msg)
			if err != nil {
				log.Printf("failed to marshal broadcast: %v", err)
				continue
			}
			for c := range h.conns {
				c.SetWriteDeadline(time.Now().Add(5 * time.Second))
				if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
					log.Printf("write to conn failed: %v", err)
					c.Close()
					delete(h.conns, c)
				}
			}
		}
	}
}

// ServeWs upgrades HTTP connection and registers it with the hub.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	// Register connection
	hub.register <- conn

	// Start reader loop for this connection (blocking until error)
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			// assume closed
			hub.unregister <- conn
			return
		}

		// Expect incoming JSON. It can be either:
		// { "type": "ready" }  => client is ready to receive history
		// OR { "user": "A", "text": "hello" } => chat message
		var incoming struct {
			Type string `json:"type"`
			User string `json:"user"`
			Text string `json:"text"`
		}
		if err := json.Unmarshal(data, &incoming); err != nil {
			log.Printf("invalid ws message: %v", err)
			continue
		}

		// If client signals ready, send history to only this connection
		if incoming.Type == "ready" {
			msgs, err := GetAllChatMessages()
			if err != nil {
				log.Printf("error fetching history for ws: %v", err)
				continue
			}
			env := wsEnvelope{Type: "history", Posts: msgs}
			b, err := json.Marshal(env)
			if err != nil {
				log.Printf("failed to marshal history envelope: %v", err)
				continue
			}
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Printf("failed to write history to conn: %v", err)
				hub.unregister <- conn
				return
			}
			log.Printf("sent history (%d messages) to client %v", len(msgs), conn.RemoteAddr())
			continue
		}

		if incoming.Text == "" {
			continue
		}

		// Save to DB and broadcast the saved message
		saved, err := SaveChatMessage(incoming.Text, incoming.User)
		if err != nil {
			log.Printf("error saving message from ws: %v", err)
			continue
		}

		env := wsEnvelope{Type: "message", Post: &saved}
		hub.broadcast <- env
	}
}

func init() {
	go hub.run()
}

package handlers

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type ChatMessage struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./chat.db")
	if err != nil {
		return err
	}

	// Create chat messages table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS chat_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT NOT NULL,
		author TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	return err
}

func SaveChatMessage(message, author string) error {
	stmt, err := db.Prepare("INSERT INTO chat_messages(message, author) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message, author)
	return err
}

func GetAllChatMessages() ([]ChatMessage, error) {
	rows, err := db.Query("SELECT id, message, author, created_at FROM chat_messages ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(&msg.ID, &msg.Message, &msg.Author, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

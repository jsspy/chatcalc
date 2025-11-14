package handlers

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbMu sync.Mutex

func InitDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Create messages table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		from_user TEXT NOT NULL,
		text TEXT NOT NULL,
		file_url TEXT,
		reply_to_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func SaveMessage(msg *Message) error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db == nil {
		log.Println("Database not initialized")
		return nil
	}

	insertSQL := `
	INSERT INTO messages (id, from_user, text, file_url, reply_to_id)
	VALUES (?, ?, ?, ?, ?);
	`

	_, err := db.Exec(insertSQL, msg.ID, msg.From, msg.Text, msg.FileURL, msg.ReplyToID)
	if err != nil {
		log.Printf("Error saving message: %v\n", err)
		return err
	}

	return nil
}

func GetMessages() ([]Message, error) {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db == nil {
		return nil, nil
	}

	selectSQL := `
	SELECT id, from_user, text, file_url, reply_to_id, created_at
	FROM messages
	ORDER BY created_at ASC;
	`

	rows, err := db.Query(selectSQL)
	if err != nil {
		log.Printf("Error querying messages: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var createdAt sql.NullString
		err := rows.Scan(&msg.ID, &msg.From, &msg.Text, &msg.FileURL, &msg.ReplyToID, &createdAt)
		if err != nil {
			log.Printf("Error scanning message: %v\n", err)
			continue
		}
		if createdAt.Valid {
			msg.CreatedAt = createdAt.String
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func CloseDB() error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		return db.Close()
	}
	return nil
}

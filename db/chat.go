package db

import (
	"database/sql"
	"github.com/cassiozareck/realchat/shared"
	"log"
	"time"
)

type ChatDB interface {
	CreateChat() (uint32, error)
	ChatExists(chatID uint32) (bool, error)
	Store(chatID uint32, msg shared.Message) error
	GetMessages(chatID uint32) ([]shared.Message, error)
}

type ChatDBImp struct {
	sql *sql.DB
}

func NewChatDBImp(sql *sql.DB) *ChatDBImp {
	return &ChatDBImp{sql}
}

// CreateChat creates a new chat and returns its ID.
func (c *ChatDBImp) CreateChat() (uint32, error) {
	var chatID uint32
	err := c.sql.QueryRow("INSERT INTO chat DEFAULT VALUES RETURNING id").Scan(&chatID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

// ChatExists checks if a chat with the given ID exists.
func (c *ChatDBImp) ChatExists(chatID uint32) (bool, error) {
	var exists bool
	err := c.sql.QueryRow("SELECT exists (SELECT 1 FROM chat WHERE id = $1)", chatID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Store will store a message in the database
func (c *ChatDBImp) Store(chatID uint32, msg shared.Message) error {

	_, err := c.sql.Exec("INSERT INTO message (chat_id, sender_id, message, timestamp) VALUES (?, ?, ?, ?)", chatID, msg.UserID, msg.Text, msg.Hour)
	if err != nil {
		return err
	}
	return nil
}

// GetMessages will get all messages from a chat
func (c *ChatDBImp) GetMessages(chatID uint32) ([]shared.Message, error) {
	rows, err := c.sql.Query("SELECT id, chat_id, sender_id, message.message, timestamp FROM message WHERE chat_id = ?", chatID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	var messages []shared.Message
	for rows.Next() {
		var msg shared.Message
		var hour string
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Text, &hour); err != nil {
			return nil, err
		}
		msg.Hour, err = time.Parse(time.RFC3339, hour)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

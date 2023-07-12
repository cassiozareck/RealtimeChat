package chat

import (
	"database/sql"
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
	"log"
	"os"
	"testing"
)

// I need a function to do the setup, create the chat and make the db connection
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

var chat *Chat
var conn *sql.DB

func setup() {
	conn = shared.ConnectToDB()
	chatDB := db.NewChatDBImp(conn)

	var err error

	chat, err = NewChat(chatDB)

	if err != nil {
		log.Fatal("Failed to create chat:", err)
	}
}

func TestChat_SendMessage(t *testing.T) {
	// Create a message
	message, err := shared.NewMessage(1, chat.id, "Hello")
	if err != nil {
		t.Error("Failed to create message:", err)
	}

	// Send the message
	err = chat.SendMessage(message)
	if err != nil {
		t.Error("Failed to send message:", err)
	}

	// Check if the message was saved in the database
	// Get the message from the database
	// Check if the message is the same as the one we created
	err = chat.UpdateMessages()
	if err != nil {
		t.Error("Failed to update messages:", err)
	}

	messages := chat.GetMessages()
	if messages[0].Text() != "Hello" || messages[0].SenderID() != 1 {
		t.Error("Message not saved in the database or not retrieved correctly")
	}
}

func shutdown() {
	err := conn.Close()
	if err != nil {
		log.Fatal("Failed to close connection:", err)
	}
}

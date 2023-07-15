package main

import (
	"database/sql"
	"encoding/json"
	"github.com/cassiozareck/realchat/chat"
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var conn *sql.DB

func main() {
	conn = shared.ConnectToDB()

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Error closing DB")
		}
	}(conn)

	http.HandleFunc("/chat", GetChat)
	http.HandleFunc("/new", NewChat)
	http.HandleFunc("/send", SendMessage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetChat(w http.ResponseWriter, r *http.Request) {
	// Parse the URL parameter
	ids, ok := r.URL.Query()["id"]

	// Ensure the id exists. If not, return an error to the user.
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
		return
	}

	// Get the first id value from the array
	id := ids[0]

	if shared.LOG {
		log.Println("Url Param 'id' is: " + string(id))
	}

	// Convert the id to uint32
	chatID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c, err := chat.GetChat(db.NewChatDBImp(conn), uint32(chatID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages, err := c.GetMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// To JSON
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewChat(w http.ResponseWriter, r *http.Request) {
	c, err := chat.NewChat(db.NewChatDBImp(conn))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(strconv.Itoa(int(c.GetID()))))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SendMessage receives a message and stores it in the database.
// example: curl -X POST -d '{"chat_id": 1, "person_id": 1, "message": "Hello"}' http://localhost:8080/send
func SendMessage(w http.ResponseWriter, r *http.Request) {
	// Check the method - it should be POST for sending message
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Error closing body")
		}
	}(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Unmarshal
	var incomingMessage shared.IncomingMessage
	err = json.Unmarshal(b, &incomingMessage)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	c, err := chat.GetChat(db.NewChatDBImp(conn), incomingMessage.ChatID)
	if err != nil {
		http.Error(w, "Error getting chat: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.SendMessage(incomingMessage)

	if err != nil {
		http.Error(w, "Error sending message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// To JSON
	_, err = w.Write([]byte(strconv.Itoa(int(c.GetID()))))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

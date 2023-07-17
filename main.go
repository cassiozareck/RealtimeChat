package main

import (
	"database/sql"
	"github.com/cassiozareck/realchat/chat"
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

	router := gin.Default()

	router.GET("/chat", GetChat)
	router.GET("/new", NewChat)
	router.POST("/send", SendMessage)

	log.Fatal(router.Run(":8080"))
}

func GetChat(c *gin.Context) {
	id := c.Query("id")

	checkAndLog("Url Param 'id' is: " + id)

	chatID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversion error: " + err.Error()})
		return
	}

	chat, err := chat.GetChat(db.NewChatDBImp(conn), uint32(chatID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	messages, err := chat.GetMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func NewChat(c *gin.Context) {
	chat, err := chat.NewChat(db.NewChatDBImp(conn))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, strconv.Itoa(int(chat.GetID())))
}

func SendMessage(c *gin.Context) {
	var incomingMessage shared.IncomingMessage

	if err := c.ShouldBindJSON(&incomingMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := chat.GetChat(db.NewChatDBImp(conn), incomingMessage.ChatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chat: " + err.Error()})
		return
	}

	err = chat.SendMessage(incomingMessage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message: " + err.Error()})
		return
	}

	c.String(http.StatusOK, strconv.Itoa(int(chat.GetID())))
}

func checkAndLog(s string) {
	if shared.LOG {
		log.Println(s)
	}
}

package main

import (
	"database/sql"
	"fmt"
	"github.com/cassiozareck/realchat/chat"
	"github.com/cassiozareck/realchat/db"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	postgres := connectToDB()

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Error closing DB")
		}
	}(postgres)

	chatDB := db.NewChatDBImp(postgres)

	talk, err := chat.GetChat(chatDB, uint32(1))

	if err != nil {
		log.Fatal(err)
	}

	messages := talk.GetMessages()
	if err != nil {
		log.Fatal(err)
	}
	for _, msg := range messages {
		fmt.Println(msg)
	}
}

func connectToDB() *sql.DB {

	// Database connection parameters
	dbHost := "postgres" // aqui tem que ser o nome da task
	dbPort := "5432"
	dbUser := "cassio"
	dbPassword := "123123"
	dbName := "realchat"

	// Create the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Fail")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected")
	return db
}

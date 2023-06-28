package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	connectToDB()
}

func connectToDB() {

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

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Error closing DB")
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected")

}

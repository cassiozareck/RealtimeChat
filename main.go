package main

import (
	"database/sql"
	"github.com/cassiozareck/realchat/shared"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	postgres := shared.ConnectToDB()

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Error closing DB")
		}
	}(postgres)

}

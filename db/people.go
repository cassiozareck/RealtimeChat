package db

import (
	"database/sql"
	"fmt"
	"github.com/cassiozareck/realchat/shared"
	"log"
)

type PeopleDB struct {
	sql *sql.DB
}

func (pdb *PeopleDB) Add(name string) error {
	query := fmt.Sprintf(`INSERT INTO person (name) VALUES ('%s')`, name)
	_, err := pdb.sql.Exec(query)
	if err != nil {
		log.Println("Error inserting data:", err)
		return err
	}
	return nil
}

func (pdb *PeopleDB) Get(id uint32) (*shared.Person, error) {
	query := fmt.Sprintf("SELECT * FROM person WHERE id = %v", id)

	rows, err := pdb.sql.Query(query)
	if err != nil {
		log.Fatal("Error querying id: ", err, "id: ", id)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal("Error closing rows")
		}
	}(rows)

	person := shared.Person{}

	rows.Next()
	err = rows.Scan(&person.ID, &person.Name)

	if err != nil {
		log.Fatal("Error scanning row:", err)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error retrieving rows:", err)
	}

	return &person, nil
}

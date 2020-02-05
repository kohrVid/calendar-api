package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect() {
	connStr := "user=calendar_api dbname=calendar_api sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

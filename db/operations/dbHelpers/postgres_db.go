package dbHelpers

import "database/sql"

func PostgresDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db

}

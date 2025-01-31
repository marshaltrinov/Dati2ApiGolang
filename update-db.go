package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Example: Insert a new record
	_, err = db.Exec("INSERT INTO locations (PostalCode, KelurahanCode) VALUES (?, ?)", "12345", "KEL001")
	if err != nil {
		log.Fatal(err)
	}
}

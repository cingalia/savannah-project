package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var (
	prefix string = "/api/v1" // API prefix
	db     *sql.DB
)

func main() {

	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/savannahdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

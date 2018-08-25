package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db = func() (d *sql.DB) {
	d, err := sql.Open("postgres", "user=lawn dbname=lawn")
	if err != nil {
		log.Fatal("database didn't open right: ", err)
	}
	return
}()

func main() {
	log.Fatal(http.ListenAndServe(":5296", nil))
}

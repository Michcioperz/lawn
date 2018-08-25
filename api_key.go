package main

import (
	"database/sql"
	"log"
)

var apiKeyValidator *sql.Stmt

func init() {
	var err error
	apiKeyValidator, err = db.Prepare("SELECT EXISTS(SELECT 1 FROM api_keys WHERE key = $1)")
	if err != nil {
		log.Fatal("statement preparation error: ", err)
	}
}

func ValidateApiKey(key string) (ok bool) {
	row := apiKeyValidator.QueryRow(key)
	err := row.Scan(&ok)
	if err != nil {
		log.Print("weird result of api key validation: ", err)
		return false
	}
	return
}

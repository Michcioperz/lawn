package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"miniflux.app/integration/nunuxkeeper"
	"net/http"
)

var nunuxInsertor *sql.Stmt

func init() {
	var err error
	nunuxInsertor, err = db.Prepare("INSERT INTO links (url, title) VALUES (?, ?)")
	if err != nil {
		log.Fatal("error preparing nunux api's statement: ", err)
	}
}

func HandleNunuxApi(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Method must be POST")
		return
	}
	user, pass, ok := req.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"You are entering somebody's lawn\"")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "You are entering somebody's lawn. Provide `api` as user and relevant API key as password.")
		return
	}
	if user != "api" || !ValidateApiKey(pass) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Get off my lawn!")
		return
	}
	var link nunuxkeeper.Document
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&link)
	if err != nil {
		log.Print("a request didn't quite parse: ", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Is that a JSON? I want to speak with JSON.")
		return
	}
	result, err := nunuxInsertor.Exec(link.Origin, link.Title)
	if err != nil {
		log.Print("link didn't insert: ", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, "Is that a link really?")
		return
	}
	log.Print("link added", result)
	w.WriteHeader(http.StatusNoContent)
	return
}

func init() {
	http.HandleFunc("/nunux/v2/documents", HandleNunuxApi)
}

package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var printAllQuery *sql.Stmt

var printAllTpl = template.Must(template.New("main").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <title>Lawn</title>
  </head>
  <body>
    <div class="container">
      <header>
        <h1>Lawn</h1>
        <h3>a primitive bookmarking solution</h3>
      </header>
      <main>
        <ul>
          {{ range . }}
          <li><a href="{{ .Url }}">{{ .Title }} <small>{{ .ParsedUrl.Hostname }}</small></a> {{ .Description }}</li>
          {{ end }}
        </ul>
      </main>
    </div>
  </body>
</html>
`))

func init() {
	var err error
	printAllQuery, err = db.Prepare("SELECT url, title, description, inserted_at FROM links")
	if err != nil {
		log.Fatal("unprepared statement for printAll: ", err)
	}
}

func PrintAll(w http.ResponseWriter, r *http.Request) {
	rows, err := printAllQuery.Query()
	if err != nil {
		log.Print("error when printing all: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal server error")
		return
	}
	defer rows.Close()
	var links []Link
	for rows.Next() {
		var link Link
		err = rows.Scan(&link.Url, &link.Title, &link.Description, &link.InsertedAt)
		if err != nil {
			log.Print("error when printing all: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal server error")
			return
		}
		link.ParsedUrl, _ = url.Parse(link.Url)
		links = append(links, link)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = printAllTpl.Execute(w, links)
}

func init() {
	http.HandleFunc("/", PrintAll)
}

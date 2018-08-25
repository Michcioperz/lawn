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
    <style>
      body {
        background-color: #4F4040;
      }
      body, a:link {
        color: #EDD3D3;
      }
      a:visited {
        color: #C2ADAD;
      }
      .container {
        width: 90%;
        max-width: 1000px;
        margin-left: auto;
        margin-right: auto;
      }
    </style>
    <link rel="alternate" type="application/atom+xml" href="/feed.atom" />
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
      <footer>
        <p>This is <a href="https://git.meekchopp.es/Michcioperz/lawn">Lawn</a>, a public bookmark list by <a href="https://michcioperz.com">Michcioperz</a>.</p>
        <p>This software is licensed AGPL v3.</p>
      </footer>
    </div>
  </body>
</html>
`))

var rssTemplate = template.Must(template.New("feed").Parse(`
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Lawn</title>
  <subtitle>a primitive bookmarking solution</subtitle>
  <link href="https://lawn.meekchopp.es/feed.atom" rel="self" />
  <link href="https://lawn.meekchopp.es/" rel="alternate" type="text/html" />
  <id>https://lawn.meekchopp.es/feed.atom</id>
  <updated>{{ if . }}{{ with $first := index . 0 }}{{ $first.InsertedAtISO }}{{ end }}{{ end }}</updated>
  {{ range . }}
  <entry>
    <title>{{ .Title }}</title>
    <link href="{{ .Url }}" />
    <id>{{ .Url }}</id>
    <updated>{{ .InsertedAtISO }}</updated>
    <summary>{{ .Description }}</summary>
  </entry>
  {{ end }}
</feed>
`))

func init() {
	var err error
	printAllQuery, err = db.Prepare("SELECT url, title, description, inserted_at FROM links ORDER BY inserted_at DESC")
	if err != nil {
		log.Fatal("unprepared statement for printAll: ", err)
	}
}

func ListAll() (links []Link, err error) {
	rows, err := printAllQuery.Query()
	defer rows.Close()
	for rows.Next() {
		var link Link
		err = rows.Scan(&link.Url, &link.Title, &link.Description, &link.InsertedAt)
		if err != nil {
			return
		}
		link.ParsedUrl, _ = url.Parse(link.Url)
		links = append(links, link)
	}
	return
}

func PrintAll(w http.ResponseWriter, r *http.Request) {
	links, err := ListAll()
	if err != nil {
		log.Print("error when printing all: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	log.Print(printAllTpl.Execute(w, links))
}

func FeedAll(w http.ResponseWriter, r *http.Request) {
	links, err := ListAll()
	if err != nil {
		log.Print("error when printing all: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `<?xml version="1.0" encoding="utf-8" ?>`)
	log.Print(rssTemplate.Execute(w, links))
}

func init() {
	http.HandleFunc("/", PrintAll)
	http.HandleFunc("/feed.atom", FeedAll)
}

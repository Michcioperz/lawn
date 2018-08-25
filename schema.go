package main

import (
	"net/url"
	"time"
)

type Link struct {
	Url         string
	Title       string
	Description string
	InsertedAt  time.Time
	ParsedUrl   *url.URL
}

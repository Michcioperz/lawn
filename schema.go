package main

import (
	"time"
)

type Link struct {
	Url         string
	Title       string
	Description string
	InsertedAt  time.Time
}

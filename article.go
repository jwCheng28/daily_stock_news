package main

import "time"

type Article struct {
	Company     string
	Title       string
	URL         string
	Source      string
	SourceURL   string
	PublishedAt time.Time
}

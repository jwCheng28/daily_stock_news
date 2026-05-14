package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func scrapeGoogleNews(company string) []Article {
	var articles []Article

	c := colly.NewCollector()

	c.OnXML("//item", func(e *colly.XMLElement) {
		title := strings.TrimSpace(e.ChildText("title"))
		link := strings.TrimSpace(e.ChildText("link"))
		source := strings.TrimSpace(e.ChildText("source"))
		sourceURL := e.ChildAttr("source", "url")
		pubDate := strings.TrimSpace(e.ChildText("pubDate"))

		if title == "" || link == "" {
			return
		}

		publishedAt, err := time.Parse(time.RFC1123, pubDate)
		if err != nil {
			log.Printf(
				"Failed to parse date %q: %v",
				pubDate,
				err,
			)
			return
		}

		articles = append(articles, Article{
            Company:     company,
			Title:       title,
			URL:         link,
			Source:      source,
			SourceURL:   sourceURL,
			PublishedAt: publishedAt,
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf(
			"Request failed: %s: %v",
			r.Request.URL,
			err,
		)
	})

	queryParam := url.QueryEscape(company)

	rssURL := fmt.Sprintf(
		"https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en",
		queryParam,
	)

	if err := c.Visit(rssURL); err != nil {
		log.Printf("Failed to visit Google News: %v", err)
	}

	return articles
}

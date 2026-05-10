package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
    "sync"
	"time"

	"github.com/gocolly/colly/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Companies []string `yaml:"companies"`
}

type Article struct {
    Company     string
	Title       string
	URL         string
	Source      string
    SourceURL   string
	PublishedAt time.Time
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

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

func deduplicateArticles(articles []Article) []Article {
	seen := make(map[string]struct{})

	result := make([]Article, 0, len(articles))

	for _, article := range articles {
		if _, exists := seen[article.URL]; exists {
			continue
		}

		seen[article.URL] = struct{}{}
		result = append(result, article)
	}

	return result
}

func main() {
	config, err := loadConfig("config/watchlist.yaml")
	if err != nil {
		log.Fatal(err)
	}

    var (
        wg sync.WaitGroup
        mu sync.Mutex
        allArticles []Article
    )

	for _, company := range config.Companies {
        wg.Add(1)

        go func(company string) {
            defer wg.Done()

            fmt.Printf("Fetching news for %s...\n", company)

            articles := scrapeGoogleNews(company)

            mu.Lock()
            allArticles = append(
                allArticles,
                articles...,
            )
            mu.Unlock()
        }(company)
	}

    wg.Wait()

	allArticles = deduplicateArticles(allArticles)

	fmt.Printf(
		"\nFound %d unique articles\n\n",
		len(allArticles),
	)

	for _, article := range allArticles {
		fmt.Printf(
			"Company: %s\n"+
				"Title: %s\n"+
				"Source: %s\n"+
				"Source URL: %s\n"+
				"Article URL: %s\n"+
				"Published: %s\n"+
				"----------------------------------------\n",
			article.Company,
			article.Title,
			article.Source,
			article.SourceURL,
			article.URL,
			article.PublishedAt.Format(time.RFC3339),
		)
	}
}

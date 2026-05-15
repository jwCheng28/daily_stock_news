package main

import (
	"fmt"
	"log"
    "sync"
	"time"
)

type SummaryResult struct {
    Company string
    Summary string
    Err     error
}

func main() {
	config, err := loadConfig("config/summarizer.yaml")
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
    allArticles = filterRecentArticles(allArticles, 24 * time.Hour)
    groupedArticles := groupByCompany(allArticles)

    ollamaWorkers := config.Ollama.Workers
    jobs := make(chan string)
    results := make(chan SummaryResult)

    var ollamaWG sync.WaitGroup

    for i := 0; i < ollamaWorkers; i++ {
        ollamaWG.Add(1)
		go func(workerID int) {
			defer ollamaWG.Done()

			for company := range jobs {
				fmt.Printf(
					"[Worker %d] Summarizing %s...\n",
					workerID,
					company,
				)

				articles := groupedArticles[company]

				sortByDate(articles)
				articles = limitArticles(
					articles,
					config.Scraping.ArticleCount,
				)

				summary, err := summarizeWithOllama(
                    config.Ollama.Model,
					company,
					articles,
				)

				results <- SummaryResult{
					Company: company,
					Summary: summary,
					Err:     err,
				}
			}
		}(i + 1)

    }

	go func() {
		for company := range groupedArticles {
			jobs <- company
		}

		close(jobs)
	}()

	go func() {
		ollamaWG.Wait()
		close(results)
	}()

    for result := range results {
		if result.Err != nil {
			log.Printf(
				"failed to summarize %s: %v",
				result.Company,
				result.Err,
			)
			continue
		}

		fmt.Printf(
			"\n================================\n"+
				"%s\n"+
				"================================\n"+
				"%s\n",
			result.Company,
			result.Summary,
		)
    }
}

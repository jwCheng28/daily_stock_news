package main

import (
	"sort"
	"time"
)

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

func filterRecentArticles(
	articles []Article,
	maxAge time.Duration,
) []Article {
	cutoff := time.Now().Add(-maxAge)

	result := make([]Article, 0, len(articles))

	for _, article := range articles {
		if article.PublishedAt.After(cutoff) {
			result = append(result, article)
		}
	}

	return result
}

func groupByCompany(
	articles []Article,
) map[string][]Article {
	grouped := make(map[string][]Article)

	for _, article := range articles {
		grouped[article.Company] = append(
			grouped[article.Company],
			article,
		)
	}

	return grouped
}

func sortByDate(articles []Article) {
	sort.Slice(
		articles,
		func(i, j int) bool {
			return articles[i].PublishedAt.After(
				articles[j].PublishedAt,
			)
		},
	)
}

func limitArticles(
	articles []Article,
	max int,
) []Article {
	if len(articles) <= max {
		return articles
	}

	return articles[:max]
}

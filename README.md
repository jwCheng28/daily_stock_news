# Financial News Aggregator & AI Summarizer

A Go-based financial news aggregator that collects company-specific news from Google News RSS, processes the articles concurrently, and uses a local Ollama LLM to generate concise financial summaries.

## Overview

The application takes a list of companies from a watchlist and:

1. Fetches & scrapes financial news of watchlist companies from Google News RSS using Colly & goroutines.
2. Deduplicates articles, filters out older articles, groups articles by company.
3. Sends each company's news to a local Ollama instance.
4. Generates an AI-powered financial summary.

## Requirements
- Go 1.XX+
- Ollama
- An Ollama-compatible local model

## Setup
### 1. Clone the repository
```sh
git clone https://github.com/jwCheng28/daily_stock_news
cd daily_stock_news
```
### 2. Install Go dependencies
```sh
go mod download
```
The project uses:
- Colly: Web/RSS scraping
- yaml.v3: YAML configuration parsing


### 3. Install Ollama
Download Ollama following their link: https://ollama.com/download

### 4. Download an Ollama Model
For example:
```sh
ollama pull llama3.2
```
### 5. Start Ollama
Run:
```bash
ollama serve
```
Ollama should be available at:
```
http://localhost:11434
```
Keep Ollama running while using the application.

### 6. Running the Application

From the project root:
```sh
go run .
```

## Configuration

The companies tracked by the application are defined in:
```
config/watchlist.yaml
```
Example:
```yaml
companies:
  - Apple
  - Microsoft
  - Nvidia
  - Tesla
  - Amazon
  - AMD
  - Meta
  - Google
```
You can add or remove companies without modifying the Go code.

## Future Improvements

Potential improvements include:

- Support for additional news sources.
- Better article relevance filtering.
- Tracking previously processed articles.
- Configurable Ollama models.
- Streaming Ollama responses.
- Scheduled/periodic news collection.
- More sophisticated financial sentiment analysis.
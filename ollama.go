package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
)

type OllamaRequest struct {
    Model string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool `json:"stream"`
}

type OllamaResponse struct {
    Response string `json:"response"`
    Done bool `json:"done"`
}

func summarizeWithOllama(
    model string,
    company string,
    articles []Article,
) (string, error) {
    var prompt strings.Builder
    prompt.WriteString(
        "You are a financial news analyst. " +
            "Summarize the following news for " +
            company +
            ". Identify: 1. Most important developments, " +
            "2. Possible market moving events, " +
            "3. Overall market sentiment. Be concise\n\n",
    )
    prompt.WriteString("News articles:\n\n")

	for i, article := range articles {
		fmt.Fprintf(
			&prompt,
			"Article %d\n"+
				"Title: %s\n"+
				"Source: %s\n",
			i+1,
			article.Title,
			article.Source,
		)
	}

	return askOllama(model, prompt.String())
}
func askOllama(model string, prompt string) (string, error) {
	requestBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf(
			"marshal request: %w",
			err,
		)
	}

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf(
			"request Ollama: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"Ollama returned HTTP %d",
			resp.StatusCode,
		)
	}

	var result OllamaResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf(
			"decode response: %w",
			err,
		)
	}

	return result.Response, nil
}

package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Companies []string `yaml:"companies"`
	Scraping ScrapingConfig `yaml:"scraping"`
	Ollama   OllamaConfig   `yaml:"ollama"`
}

type ScrapingConfig struct {
	ArticleCount int `yaml:"article_count"`
}

type OllamaConfig struct {
	Model   string `yaml:"model"`
	Workers int    `yaml:"workers"`
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if len(config.Companies) == 0 {
		return Config{}, fmt.Errorf(
			"watchlist contains no companies",
		)
	}


	if config.Scraping.ArticleCount <= 0 {
		config.Scraping.ArticleCount = 5
	}

	if config.Ollama.Model == "" {
		config.Ollama.Model = "llama3.2"
	}

	if config.Ollama.Workers <= 0 {
		config.Ollama.Workers = 3
	}


	return config, nil
}

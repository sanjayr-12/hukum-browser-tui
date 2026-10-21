package test

import (
	"context"
	"os"
	"testing"

	"hukum/internal/config"
	"hukum/internal/tinyfish"
)

func init() {
	config.LoadEnvFile(".env")
	config.LoadEnvFile("../.env")
}

func TestClientFetch(t *testing.T) {
	apiKey := os.Getenv("TINYFISH_API_KEY")
	if apiKey == "" {
		t.Skip("TINYFISH_API_KEY not set, skipping integration test")
	}

	client := tinyfish.NewClient(apiKey)
	res, err := client.Fetch(context.Background(), "https://news.ycombinator.com")
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}

	if res.Title == "" {
		t.Errorf("expected non-empty title")
	}
	if len(res.Text) == 0 {
		t.Errorf("expected non-empty text")
	}
}

func TestClientSearch(t *testing.T) {
	apiKey := os.Getenv("TINYFISH_API_KEY")
	if apiKey == "" {
		t.Skip("TINYFISH_API_KEY not set, skipping integration test")
	}

	client := tinyfish.NewClient(apiKey)
	results, err := client.Search(context.Background(), "hacker news")
	if err != nil {
		t.Fatalf("search error: %v", err)
	}

	if len(results) == 0 {
		t.Errorf("expected at least 1 search result")
	}
}

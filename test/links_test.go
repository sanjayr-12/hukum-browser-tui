package test

import (
	"testing"

	"hukum/internal/ui"
)

func TestExtractLinks(t *testing.T) {
	md := `
# Welcome to Hacker News
Check out [PHP 8.2](https://github.com/php/php-src) and [Show HN: Browser](https://news.ycombinator.com/item?id=12345).
Also visit [Internal Page](/about).
`
	raw := []string{
		"https://news.ycombinator.com/newest",
		"https://news.ycombinator.com/front",
	}

	links := ui.ExtractLinks(md, raw, "https://news.ycombinator.com")
	if len(links) < 3 {
		t.Fatalf("expected at least 3 links, got %d", len(links))
	}

	if links[0].Title != "PHP 8.2" {
		t.Errorf("expected first link title 'PHP 8.2', got '%s'", links[0].Title)
	}
	if links[0].URL != "https://github.com/php/php-src" {
		t.Errorf("expected first link URL 'https://github.com/php/php-src', got '%s'", links[0].URL)
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://news.ycombinator.com", true},
		{"http://example.com/test", true},
		{"golang.org", true},
		{"github.com/charmbracelet/bubbletea", true},
		{"what is the best programming language", false},
		{"latest news on AI", false},
	}

	for _, tt := range tests {
		got := ui.IsURL(tt.input)
		if got != tt.expected {
			t.Errorf("IsURL(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

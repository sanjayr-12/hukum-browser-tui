package test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/glamour"
	"hukum/internal/tinyfish"
	"hukum/internal/ui"
)

func TestGlamourRendering(t *testing.T) {
	md := "# Test Heading\n\nThis is a **bold** paragraph with a [link](https://example.com).\n\n* Item 1\n* Item 2\n\n```go\nfmt.Println(\"Hello world\")\n```"
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create glamour renderer: %v", err)
	}

	out, err := renderer.Render(md)
	if err != nil {
		t.Fatalf("failed to render markdown: %v", err)
	}

	if !strings.Contains(out, "Test") || !strings.Contains(out, "Heading") {
		t.Errorf("rendered output missing heading words")
	}
}

func TestFullPageSearchResultsRendering(t *testing.T) {
	modal := ui.NewModal()
	modal.OpenSearch()
	modal.SetSearchResults([]tinyfish.SearchResultItem{
		{
			Title:   "React Native · Learn once, write anywhere",
			URL:     "https://reactnative.dev/",
			Snippet: "React Native lets you create truly native apps and doesn't compromise...",
		},
		{
			Title:   "React Native - Wikipedia",
			URL:     "https://en.wikipedia.org/wiki/React_Native",
			Snippet: "React Native is an open-source UI software framework...",
		},
	})

	output := modal.RenderFullPage(100, 24)
	if len(output) == 0 {
		t.Fatalf("expected non-empty rendered output")
	}

	if !strings.Contains(output, "React Native") {
		t.Errorf("output does not contain result title")
	}
	if !strings.Contains(output, "https://reactnative.dev/") {
		t.Errorf("output does not contain result URL")
	}
}

package test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/glamour"
	"hukum/internal/tinyfish"
	"hukum/internal/ui"
)

func TestOneDarkGlamourRendering(t *testing.T) {
	md := "# Test Heading\n\nThis is a **bold** paragraph with a [link](https://example.com).\n\n* Item 1\n* Item 2\n\n```go\nfmt.Println(\"Hello world\")\n```"
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(ui.OneDarkStyleConfig),
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

func TestAdStripping(t *testing.T) {
	rawWithAd := `# Clean Article Title

This is legitimate article text that should be preserved.

comcastbusiness •
Ad
Get the most comprehensive cybersecurity for small business at the best price + ask about a $500 prepaid card.
Learn More
comcastbusiness.com

This is the next paragraph after the ad.`

	cleaned := ui.NormalizeMarkdown(rawWithAd)

	if strings.Contains(cleaned, "comcastbusiness") {
		t.Errorf("ad block was not stripped: %s", cleaned)
	}
	if strings.Contains(cleaned, "cybersecurity for small business") {
		t.Errorf("ad copy was not stripped: %s", cleaned)
	}
	if !strings.Contains(cleaned, "Clean Article Title") {
		t.Errorf("expected article title to be preserved")
	}
	if !strings.Contains(cleaned, "This is the next paragraph") {
		t.Errorf("expected following paragraph to be preserved")
	}
}

func TestCodeBlockNormalization(t *testing.T) {
	rawGFGCode := `Java

` + "```` ```" + `
public class Main {
    public static void main(String[] args) {
        int a = 10;
        int b = 20;
        int sum = a + b;
        System.out.println("Sum = " + sum);
    }
}
` + "``` ````"

	cleaned := ui.NormalizeMarkdown(rawGFGCode)
	t.Logf("Cleaned markdown:\n%s\n", cleaned)

	if strings.Contains(cleaned, "````") {
		t.Errorf("quadruple backticks were not cleaned: %s", cleaned)
	}
	if !strings.Contains(cleaned, "```Java") && !strings.Contains(cleaned, "```java") {
		t.Errorf("language tag was not attached to code block: %s", cleaned)
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(ui.OneDarkStyleConfig),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	out, err := renderer.Render(cleaned)
	if err != nil {
		t.Fatalf("failed to render normalized code: %v", err)
	}

	t.Logf("Rendered code block:\n%s\n", out)
	if !strings.Contains(out, "Main") || !strings.Contains(out, "println") {
		t.Errorf("rendered code missing class definition")
	}
}

func TestNormalizeMarkdownHackerNews(t *testing.T) {
	rawHN := `Hacker News
new
|
past
1.
PHP 8.2.35
(
github.com/php
)
827 points
by
tuyuri
4 hours ago
|
hide
|
467 comments
2.
Today's PHP security release
(
thephp.foundation
)
288 points
by
ejohnso49
2 hours ago
|
hide
|
147 comments`

	normalized := ui.NormalizeMarkdown(rawHN)

	if !strings.Contains(normalized, "PHP 8.2.35") {
		t.Errorf("expected normalized item 1, got: %s", normalized)
	}
	if !strings.Contains(normalized, "(github.com/php)") {
		t.Errorf("expected normalized domain in parentheses, got: %s", normalized)
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(ui.OneDarkStyleConfig),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	out, err := renderer.Render(normalized)
	if err != nil {
		t.Fatalf("failed to render normalized markdown: %v", err)
	}

	if !strings.Contains(out, "PHP 8.2.35") || !strings.Contains(out, "Today's PHP security release") {
		t.Errorf("glamour output does not contain story titles")
	}
}

func TestCenteredSearchPromptDialog(t *testing.T) {
	modal := ui.NewModal()
	modal.OpenSearch()

	if !modal.IsInputOnly() {
		t.Fatalf("expected IsInputOnly to be true for fresh search modal")
	}

	output := modal.RenderDialog(80, 20)
	if len(output) == 0 {
		t.Fatalf("expected non-empty rendered output")
	}

	if !strings.Contains(output, "Search or URL") {
		t.Errorf("output missing modal title")
	}
}

func TestSearchResultsDialog(t *testing.T) {
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

	if modal.IsInputOnly() {
		t.Fatalf("expected IsInputOnly to be false after results loaded")
	}

	output := modal.RenderDialog(100, 24)
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

func TestOverlayCenter(t *testing.T) {
	bg := strings.Repeat("background background background\n", 10)
	fg := "┌────────┐\n│ dialog │\n└────────┘"

	combined := ui.OverlayCenter(bg, fg, 40, 10)
	if len(combined) == 0 {
		t.Fatalf("expected non-empty combined output")
	}

	if !strings.Contains(combined, "dialog") {
		t.Errorf("expected combined output to contain dialog text")
	}
	if !strings.Contains(combined, "background") {
		t.Errorf("expected background to remain visible around dialog")
	}
}

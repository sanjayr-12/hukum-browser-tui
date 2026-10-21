package ui

import (
	"net/url"
	"regexp"
	"strings"
)

type PageLink struct {
	Title string
	URL   string
}

var mdLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^\s\)\"]+)\)`)

func ExtractLinks(mdText string, rawLinks []string, baseURL string) []PageLink {
	var results []PageLink
	seen := make(map[string]bool)

	matches := mdLinkRegex.FindAllStringSubmatch(mdText, -1)
	for _, m := range matches {
		if len(m) >= 3 {
			title := strings.TrimSpace(m[1])
			linkURL := strings.TrimSpace(m[2])
			if title != "" && linkURL != "" && !seen[linkURL] {
				seen[linkURL] = true
				results = append(results, PageLink{
					Title: cleanLinkTitle(title),
					URL:   linkURL,
				})
			}
		}
	}

	var parsedBase *url.URL
	if baseURL != "" {
		parsedBase, _ = url.Parse(baseURL)
	}

	for _, raw := range rawLinks {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		fullURL := raw
		if parsedBase != nil && !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
			if rel, err := url.Parse(raw); err == nil {
				fullURL = parsedBase.ResolveReference(rel).String()
			}
		}

		if !strings.HasPrefix(fullURL, "http://") && !strings.HasPrefix(fullURL, "https://") {
			continue
		}

		if !seen[fullURL] {
			seen[fullURL] = true
			results = append(results, PageLink{
				Title: formatURLAsTitle(fullURL),
				URL:   fullURL,
			})
		}
	}

	return results
}

func cleanLinkTitle(t string) string {
	t = strings.ReplaceAll(t, "\n", " ")
	t = strings.ReplaceAll(t, "\r", "")
	t = strings.Join(strings.Fields(t), " ")
	if len(t) > 60 {
		return t[:57] + "..."
	}
	return t
}

func formatURLAsTitle(u string) string {
	if parsed, err := url.Parse(u); err == nil {
		path := strings.TrimPrefix(parsed.Path, "/")
		if path != "" {
			return parsed.Host + "/" + path
		}
		return parsed.Host
	}
	return u
}

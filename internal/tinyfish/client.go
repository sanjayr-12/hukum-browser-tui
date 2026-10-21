package tinyfish

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultSearchEndpoint = "https://api.search.tinyfish.ai"
	DefaultFetchEndpoint  = "https://api.fetch.tinyfish.ai"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type SearchResultItem struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type SearchResponse struct {
	Results []SearchResultItem `json:"results"`
}

type FetchResultItem struct {
	URL         string   `json:"url"`
	FinalURL    string   `json:"final_url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Format      string   `json:"format"`
	Text        string   `json:"text"`
	Links       []string `json:"links,omitempty"`
	ImageLinks  []string `json:"image_links,omitempty"`
}

type FetchErrorItem struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

type FetchResponse struct {
	Results []FetchResultItem `json:"results"`
	Errors  []FetchErrorItem  `json:"errors"`
}

type FetchRequestPayload struct {
	URLs       []string `json:"urls"`
	Format     string   `json:"format"`
	Links      bool     `json:"links"`
	ImageLinks bool     `json:"image_links"`
}

func NewClient(apiKey string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("TINYFISH_API_KEY")
	}
	return &Client{
		apiKey: strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

func (c *Client) HasAPIKey() bool {
	return c.apiKey != ""
}

func (c *Client) SetAPIKey(key string) {
	c.apiKey = strings.TrimSpace(key)
}

func (c *Client) Search(ctx context.Context, query string) ([]SearchResultItem, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TINYFISH_API_KEY is not set. Please set the environment variable or pass an API key")
	}

	endpoint, err := url.Parse(DefaultSearchEndpoint)
	if err != nil {
		return nil, err
	}

	q := endpoint.Query()
	q.Set("query", query)
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading search response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search API returned HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(bodyBytes, &searchResp); err != nil {
		return nil, fmt.Errorf("parsing search response failed: %w", err)
	}

	return searchResp.Results, nil
}

func (c *Client) Fetch(ctx context.Context, targetURL string) (*FetchResultItem, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TINYFISH_API_KEY is not set. Please set the environment variable or pass an API key")
	}

	targetURL = strings.TrimSpace(targetURL)
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	payload := FetchRequestPayload{
		URLs:       []string{targetURL},
		Format:     "markdown",
		Links:      true,
		ImageLinks: true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, DefaultFetchEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading fetch response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch API returned HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var fetchResp FetchResponse
	if err := json.Unmarshal(bodyBytes, &fetchResp); err != nil {
		return nil, fmt.Errorf("parsing fetch response failed: %w", err)
	}

	if len(fetchResp.Results) == 0 {
		if len(fetchResp.Errors) > 0 {
			return nil, fmt.Errorf("fetch error: %s", fetchResp.Errors[0].Error)
		}
		return nil, fmt.Errorf("no content returned for %s", targetURL)
	}

	return &fetchResp.Results[0], nil
}

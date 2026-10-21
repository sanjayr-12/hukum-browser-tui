package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	APIKey      string
	InitialURL  string
	SearchQuery string
}

func Load() *Config {
	return LoadArgs(os.Args[1:])
}

func LoadArgs(args []string) *Config {
	LoadEnvFile(".env")

	var (
		apiKeyFlag  string
		initialURL  string
		searchQuery string
	)

	fs := flag.NewFlagSet("hukum", flag.ContinueOnError)
	fs.StringVar(&apiKeyFlag, "api-key", "", "TinyFish API Key")
	fs.StringVar(&initialURL, "url", "", "Initial URL to open")
	fs.StringVar(&searchQuery, "search", "", "Initial search query")
	if err := fs.Parse(args); err != nil {
		if strings.Contains(err.Error(), "help requested") {
			os.Exit(0)
		}
	}

	if initialURL == "" && fs.NArg() > 0 {
		initialURL = fs.Arg(0)
	}

	apiKey := apiKeyFlag
	if apiKey == "" {
		apiKey = os.Getenv("TINYFISH_API_KEY")
	}

	return &Config{
		APIKey:      strings.TrimSpace(apiKey),
		InitialURL:  strings.TrimSpace(initialURL),
		SearchQuery: strings.TrimSpace(searchQuery),
	}
}

func LoadEnvFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

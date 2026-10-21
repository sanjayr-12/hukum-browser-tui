package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"hukum/internal/config"
	"hukum/internal/tinyfish"
	"hukum/internal/ui"
)

func main() {
	cfg := config.Load()
	client := tinyfish.NewClient(cfg.APIKey)

	if !client.HasAPIKey() {
		fmt.Println("Warning: TINYFISH_API_KEY is not set.")
		fmt.Println("Set it in a .env file or run with:")
		fmt.Println("  export TINYFISH_API_KEY=\"your_key\"")
	}

	model := ui.NewBrowserModel(client, cfg.InitialURL, cfg.SearchQuery)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running browser: %v\n", err)
		os.Exit(1)
	}
}

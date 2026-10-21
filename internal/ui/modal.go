package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"hukum/internal/tinyfish"
)

type ModalType int

const (
	ModalNone ModalType = iota
	ModalSearch
	ModalLinks
)

type ModalItem struct {
	Title   string
	URL     string
	Snippet string
}

type ModalModel struct {
	Type        ModalType
	Input       textinput.Model
	Items       []ModalItem
	Filtered    []ModalItem
	Cursor      int
	IsSearching bool
	SearchErr   string
	LastQuery   string
}

func NewModal() ModalModel {
	ti := textinput.New()
	ti.Prompt = "Search / URL: "
	ti.PromptStyle = SearchPromptStyle
	ti.TextStyle = NormalRowStyle
	ti.Focus()

	return ModalModel{
		Type:  ModalNone,
		Input: ti,
	}
}

func (m *ModalModel) OpenSearch() {
	m.Type = ModalSearch
	m.Input.Reset()
	m.Input.Prompt = "Search or URL: "
	m.Input.Placeholder = "Enter query (e.g. latest tech news) or URL (e.g. news.ycombinator.com)..."
	m.Input.Focus()
	m.Items = nil
	m.Filtered = nil
	m.Cursor = 0
	m.IsSearching = false
	m.SearchErr = ""
	m.LastQuery = ""
}

func (m *ModalModel) OpenLinks(links []PageLink) {
	m.Type = ModalLinks
	m.Input.Reset()
	m.Input.Prompt = "Filter Links: "
	m.Input.Placeholder = "Type keywords to filter page links..."
	m.Input.Focus()

	items := make([]ModalItem, 0, len(links))
	for _, l := range links {
		items = append(items, ModalItem{
			Title: l.Title,
			URL:   l.URL,
		})
	}
	m.Items = items
	m.Filtered = items
	m.Cursor = 0
	m.IsSearching = false
	m.SearchErr = ""
}

func (m *ModalModel) Close() {
	m.Type = ModalNone
	m.Input.Blur()
}

func (m *ModalModel) Filter() {
	if m.Type != ModalLinks {
		return
	}
	query := strings.ToLower(strings.TrimSpace(m.Input.Value()))
	if query == "" {
		m.Filtered = m.Items
		m.Cursor = 0
		return
	}

	var filtered []ModalItem
	for _, item := range m.Items {
		if strings.Contains(strings.ToLower(item.Title), query) ||
			strings.Contains(strings.ToLower(item.URL), query) {
			filtered = append(filtered, item)
		}
	}
	m.Filtered = filtered
	if m.Cursor >= len(filtered) {
		m.Cursor = max(0, len(filtered)-1)
	}
}

func (m *ModalModel) SetSearchResults(results []tinyfish.SearchResultItem) {
	items := make([]ModalItem, 0, len(results))
	for _, r := range results {
		items = append(items, ModalItem{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Snippet,
		})
	}
	m.Items = items
	m.Filtered = items
	m.Cursor = 0
	m.IsSearching = false
	m.SearchErr = ""
	m.LastQuery = strings.TrimSpace(m.Input.Value())
}

func (m *ModalModel) Next() {
	if len(m.Filtered) > 0 {
		m.Cursor = (m.Cursor + 1) % len(m.Filtered)
	}
}

func (m *ModalModel) Prev() {
	if len(m.Filtered) > 0 {
		m.Cursor = (m.Cursor - 1 + len(m.Filtered)) % len(m.Filtered)
	}
}

func (m *ModalModel) SelectedItem() *ModalItem {
	if len(m.Filtered) > 0 && m.Cursor >= 0 && m.Cursor < len(m.Filtered) {
		return &m.Filtered[m.Cursor]
	}
	return nil
}

// RenderFullPage renders the search or link picker spanning the full terminal width and height
func (m *ModalModel) RenderFullPage(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	// 1. Search / Input Bar (Full Width)
	counterText := ""
	if len(m.Filtered) > 0 {
		counterText = fmt.Sprintf("[%d / %d]", m.Cursor+1, len(m.Filtered))
	} else if m.IsSearching {
		counterText = "[searching...]"
	}
	counterBadge := CounterBadgeStyle.Render(counterText)

	inputWidth := width - lipgloss.Width(counterBadge) - lipgloss.Width(m.Input.Prompt) - 4
	if inputWidth < 20 {
		inputWidth = 20
	}
	m.Input.Width = inputWidth

	inputRowContent := m.Input.View()
	spacing := width - lipgloss.Width(inputRowContent) - lipgloss.Width(counterBadge) - 2
	if spacing < 1 {
		spacing = 1
	}
	searchBarLine := SearchContainerStyle.Width(width).Render(
		inputRowContent + strings.Repeat(" ", spacing) + counterBadge,
	)

	// Separator line
	dividerLine := DividerStyle.Render(strings.Repeat("─", width))

	// Available height for results list
	// 1 line for search bar, 1 line for divider, 1 line for bottom status bar
	listHeight := height - 3
	if listHeight < 4 {
		listHeight = 4
	}

	var lines []string

	if m.IsSearching {
		lines = append(lines, "")
		lines = append(lines, "  "+SpinnerStyle.Render("Searching TinyFish web index..."))
	} else if m.SearchErr != "" {
		lines = append(lines, "")
		lines = append(lines, "  "+ErrorStyle.Render("Search Error: "+m.SearchErr))
	} else if len(m.Filtered) == 0 {
		lines = append(lines, "")
		if m.Type == ModalSearch {
			if strings.TrimSpace(m.Input.Value()) == "" {
				lines = append(lines, "  "+TextMutedStyle.Render("Type a search query or a website URL (e.g. news.ycombinator.com), then press Enter."))
			} else {
				lines = append(lines, "  "+TextMutedStyle.Render("Press Enter to search TinyFish or navigate to URL."))
			}
		} else {
			lines = append(lines, "  "+TextMutedStyle.Render("No matching links found on this page."))
		}
	} else {
		// Calculate lines per item: 3 lines for search results, 2 lines for links
		linesPerItem := 3
		if m.Type == ModalLinks {
			linesPerItem = 2
		}

		itemsPerPage := listHeight / linesPerItem
		if itemsPerPage < 1 {
			itemsPerPage = 1
		}

		startIdx := 0
		if m.Cursor >= itemsPerPage {
			startIdx = m.Cursor - itemsPerPage + 1
		}
		endIdx := startIdx + itemsPerPage
		if endIdx > len(m.Filtered) {
			endIdx = len(m.Filtered)
		}

		for i := startIdx; i < endIdx; i++ {
			item := m.Filtered[i]
			isSel := i == m.Cursor

			titleStr := item.Title
			if titleStr == "" {
				titleStr = item.URL
			}
			maxLen := width - 6
			if maxLen > 0 && len(titleStr) > maxLen {
				titleStr = titleStr[:maxLen-3] + "..."
			}

			var titleLine string
			if isSel {
				titleLine = SelectedRowStyle.Width(width).Render(" > " + titleStr)
			} else {
				titleLine = NormalRowStyle.Render("   " + titleStr)
			}
			lines = append(lines, titleLine)

			urlStr := item.URL
			if maxLen > 0 && len(urlStr) > maxLen {
				urlStr = urlStr[:maxLen-3] + "..."
			}
			lines = append(lines, "     "+ResultURLStyle.Render(urlStr))

			if m.Type == ModalSearch {
				snippetStr := strings.TrimSpace(strings.ReplaceAll(item.Snippet, "\n", " "))
				if maxLen > 0 && len(snippetStr) > maxLen {
					snippetStr = snippetStr[:maxLen-3] + "..."
				}
				if snippetStr != "" {
					lines = append(lines, "     "+ResultSnippetStyle.Render(snippetStr))
				} else {
					lines = append(lines, "")
				}
			}
		}
	}

	for len(lines) < listHeight {
		lines = append(lines, "")
	}
	if len(lines) > listHeight {
		lines = lines[:listHeight]
	}
	contentArea := strings.Join(lines, "\n")

	// Bottom Status Bar
	statusBar := StatusBarStyle.Width(width).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		StatusKeyStyle.Render(" Enter "), StatusValStyle.Render("Open/Search  "),
		StatusKeyStyle.Render("Ctrl+N/P "), StatusValStyle.Render("Navigate  "),
		StatusKeyStyle.Render("Tab "), StatusValStyle.Render("Next  "),
		StatusKeyStyle.Render("Esc "), StatusValStyle.Render("Back to page  "),
		StatusKeyStyle.Render("Ctrl+C "), StatusValStyle.Render("Quit"),
	))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		searchBarLine,
		dividerLine,
		contentArea,
		statusBar,
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

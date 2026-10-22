package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
	ti.Prompt = "> "
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
	m.Input.Prompt = "> "
	m.Input.Placeholder = "Search query or URL (e.g. news.ycombinator.com)..."
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
	m.Input.Prompt = "> "
	m.Input.Placeholder = "Filter page links..."
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

// IsInputOnly returns true when the search modal should show just the small
// centered prompt (no results loaded yet, not currently searching).
func (m *ModalModel) IsInputOnly() bool {
	return m.Type == ModalSearch && len(m.Items) == 0 && !m.IsSearching && m.SearchErr == ""
}

// RenderDialog renders the modal as a floating dialog box.
func (m *ModalModel) RenderDialog(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	if m.IsInputOnly() {
		return m.renderInputPrompt(width)
	}
	return m.renderResultsDialog(width, height)
}

func (m *ModalModel) renderInputPrompt(screenWidth int) string {
	promptWidth := int(float64(screenWidth) * 0.55)
	if promptWidth < 40 {
		promptWidth = screenWidth - 4
	}
	if promptWidth > 72 {
		promptWidth = 72
	}

	boxWidth := promptWidth - 2
	innerWidth := boxWidth - 2
	if innerWidth < 10 {
		innerWidth = 10
	}

	titleBar := m.renderTitleBar(" Search or URL ", innerWidth)

	m.Input.Width = innerWidth - lipgloss.Width(m.Input.Prompt) - 1
	if m.Input.Width < 15 {
		m.Input.Width = 15
	}
	inputRow := m.Input.View()
	inputW := ansi.StringWidth(inputRow)
	if inputW < innerWidth {
		inputRow += strings.Repeat(" ", innerWidth-inputW)
	}

	divider := ModalDividerStyle.Render(strings.Repeat("─", innerWidth))
	footerHints := ModalFooterStyle.Render("Enter Search  Esc Cancel")
	footerW := ansi.StringWidth(footerHints)
	if footerW < innerWidth {
		footerHints += strings.Repeat(" ", innerWidth-footerW)
	}

	modalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleBar,
		inputRow,
		divider,
		footerHints,
	)

	return ModalBoxStyle.Width(boxWidth).Render(modalContent)
}

func (m *ModalModel) renderResultsDialog(screenWidth, screenHeight int) string {
	modalWidth := int(float64(screenWidth) * 0.80)
	if modalWidth < 60 {
		modalWidth = screenWidth - 4
	}
	if modalWidth > 110 {
		modalWidth = 110
	}

	modalHeight := int(float64(screenHeight) * 0.70)
	if modalHeight < 14 {
		modalHeight = screenHeight - 2
	}
	if modalHeight > 30 {
		modalHeight = 30
	}

	boxWidth := modalWidth - 2
	boxHeight := modalHeight - 2
	innerWidth := boxWidth - 2
	if innerWidth < 10 {
		innerWidth = 10
	}

	var titleText string
	if m.Type == ModalSearch {
		if m.IsSearching {
			titleText = " Searching... "
		} else {
			titleText = " Search Results "
		}
	} else {
		titleText = " Page Links "
	}
	titleBar := m.renderTitleBar(titleText, innerWidth)

	counterText := ""
	if len(m.Filtered) > 0 {
		counterText = fmt.Sprintf("[%d / %d]", m.Cursor+1, len(m.Filtered))
	} else if m.IsSearching {
		counterText = "[searching...]"
	}
	counterBadge := CounterBadgeStyle.Render(counterText)

	inputWidth := innerWidth - lipgloss.Width(counterBadge) - lipgloss.Width(m.Input.Prompt) - 2
	if inputWidth < 15 {
		inputWidth = 15
	}
	m.Input.Width = inputWidth

	inputContent := m.Input.View()
	promptSpacing := innerWidth - ansi.StringWidth(inputContent) - ansi.StringWidth(counterBadge)
	if promptSpacing < 1 {
		promptSpacing = 1
	}
	inputRow := inputContent + strings.Repeat(" ", promptSpacing) + counterBadge

	divider := ModalDividerStyle.Render(strings.Repeat("─", innerWidth))

	listHeight := boxHeight - 5
	if listHeight < 3 {
		listHeight = 3
	}

	var lines []string
	if m.IsSearching {
		lines = append(lines, "")
		lines = append(lines, "  "+SpinnerStyle.Render("Searching TinyFish web index..."))
	} else if m.SearchErr != "" {
		lines = append(lines, "")
		lines = append(lines, "  "+ErrorStyle.Render("Error: "+m.SearchErr))
	} else if len(m.Filtered) == 0 {
		lines = append(lines, "")
		if m.Type == ModalSearch {
			lines = append(lines, "  "+TextMutedStyle.Render("No results found. Edit your query and press Enter."))
		} else {
			lines = append(lines, "  "+TextMutedStyle.Render("No matching links."))
		}
	} else {
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

		maxLen := innerWidth - 6
		for i := startIdx; i < endIdx; i++ {
			item := m.Filtered[i]
			isSel := i == m.Cursor

			titleStr := item.Title
			if titleStr == "" {
				titleStr = item.URL
			}
			if maxLen > 0 && len(titleStr) > maxLen {
				titleStr = titleStr[:maxLen-3] + "..."
			}

			if isSel {
				lines = append(lines, SelectedRowStyle.Width(innerWidth).Render(" > "+titleStr))
			} else {
				lines = append(lines, NormalRowStyle.Render("   "+titleStr))
			}

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

	footerHints := ModalFooterStyle.Render("Enter Open  Ctrl+N/P Navigate  Tab Next  Esc Close")
	footerW := ansi.StringWidth(footerHints)
	if footerW < innerWidth {
		footerHints += strings.Repeat(" ", innerWidth-footerW)
	}

	modalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleBar,
		inputRow,
		divider,
		contentArea,
		divider,
		footerHints,
	)

	return ModalBoxStyle.Width(boxWidth).Height(boxHeight).Render(modalContent)
}

func (m *ModalModel) renderTitleBar(title string, innerWidth int) string {
	titleStyled := ModalTitleStyle.Render(title)
	dashesCount := innerWidth - lipgloss.Width(titleStyled) - 2
	if dashesCount < 0 {
		dashesCount = 0
	}
	return ModalDividerStyle.Render("──") + titleStyled + ModalDividerStyle.Render(strings.Repeat("─", dashesCount))
}

// OverlayCenter places the overlay string centered on top of the background
// string, replacing only the area covered by the overlay. The surrounding
// background content remains visible.
func OverlayCenter(bg, fg string, width, height int) string {
	if width <= 0 || height <= 0 {
		return bg
	}

	bgLines := strings.Split(bg, "\n")
	fgLines := strings.Split(fg, "\n")

	for len(bgLines) < height {
		bgLines = append(bgLines, strings.Repeat(" ", width))
	}
	if len(bgLines) > height {
		bgLines = bgLines[:height]
	}

	fgHeight := len(fgLines)
	fgWidth := 0
	for _, line := range fgLines {
		if w := ansi.StringWidth(line); w > fgWidth {
			fgWidth = w
		}
	}

	startY := (height - fgHeight) / 2
	startX := (width - fgWidth) / 2
	if startY < 0 {
		startY = 0
	}
	if startX < 0 {
		startX = 0
	}

	for i, fgLine := range fgLines {
		y := startY + i
		if y >= len(bgLines) {
			break
		}

		bgLine := bgLines[y]
		bgVisualWidth := ansi.StringWidth(bgLine)
		if bgVisualWidth < width {
			bgLine = bgLine + strings.Repeat(" ", width-bgVisualWidth)
		}

		left := ansi.Cut(bgLine, 0, startX)
		leftWidth := ansi.StringWidth(left)
		if leftWidth < startX {
			left += strings.Repeat(" ", startX-leftWidth)
		}

		fgLineWidth := ansi.StringWidth(fgLine)
		if fgLineWidth < fgWidth {
			fgLine += strings.Repeat(" ", fgWidth-fgLineWidth)
		}

		rightStart := startX + fgWidth
		var right string
		if rightStart < width {
			right = ansi.Cut(bgLine, rightStart, width)
			rightWidth := ansi.StringWidth(right)
			if rightWidth < width-rightStart {
				right += strings.Repeat(" ", (width-rightStart)-rightWidth)
			}
		}

		bgLines[y] = left + fgLine + right
	}

	return strings.Join(bgLines, "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

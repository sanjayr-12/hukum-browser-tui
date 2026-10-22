package ui

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"hukum/internal/tinyfish"
)

type fetchSuccessMsg struct {
	item    *tinyfish.FetchResultItem
	latency time.Duration
}

type fetchErrorMsg struct {
	err error
}

type searchSuccessMsg struct {
	results []tinyfish.SearchResultItem
}

type searchErrorMsg struct {
	err error
}

type BrowserModel struct {
	client      *tinyfish.Client
	viewport    viewport.Model
	modal       ModalModel
	spinner     spinner.Model
	isLoading   bool
	loadingMsg  string
	currentURL  string
	pageTitle   string
	rawMarkdown string
	renderedMD  string
	pageLinks   []PageLink
	history     []string
	err         error
	width       int
	height      int
	ready       bool
}

func NewBrowserModel(client *tinyfish.Client, initialURL, initialSearch string) BrowserModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = SpinnerStyle

	m := BrowserModel{
		client:  client,
		modal:   NewModal(),
		spinner: sp,
	}

	if initialURL != "" {
		m.currentURL = initialURL
	} else if initialSearch != "" {
		m.modal.OpenSearch()
		m.modal.Input.SetValue(initialSearch)
		m.modal.IsSearching = true
	} else {
		m.modal.OpenSearch()
	}

	return m
}

func (m BrowserModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.spinner.Tick)
	if m.currentURL != "" {
		cmds = append(cmds, m.fetchURL(m.currentURL))
	} else if m.modal.IsSearching && m.modal.Input.Value() != "" {
		cmds = append(cmds, m.searchQuery(m.modal.Input.Value()))
	}
	return tea.Batch(cmds...)
}

func (m BrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 1
		footerHeight := 1
		verticalMargin := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMargin)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargin
		}

		if m.rawMarkdown != "" {
			m.renderMarkdown()
			m.viewport.SetContent(m.renderedMD)
		}

	case spinner.TickMsg:
		if m.isLoading || m.modal.IsSearching {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			cmds = append(cmds, spinCmd)
		}

	case fetchSuccessMsg:
		m.isLoading = false
		m.err = nil
		m.currentURL = msg.item.FinalURL
		if m.currentURL == "" {
			m.currentURL = msg.item.URL
		}
		m.pageTitle = msg.item.Title
		m.rawMarkdown = msg.item.Text
		m.pageLinks = ExtractLinks(m.rawMarkdown, msg.item.Links, m.currentURL)

		m.renderMarkdown()
		m.viewport.SetContent(m.renderedMD)
		m.viewport.GotoTop()

	case fetchErrorMsg:
		m.isLoading = false
		m.err = msg.err
		errorMD := fmt.Sprintf("# Error Loading Page\n\n**%s**\n\nPress `Ctrl+O` to search or enter another URL, or `b` to go back.", msg.err.Error())
		m.renderedMD = errorMD
		m.viewport.SetContent(ErrorStyle.Render(errorMD))

	case searchSuccessMsg:
		m.modal.IsSearching = false
		m.modal.SetSearchResults(msg.results)

	case searchErrorMsg:
		m.modal.IsSearching = false
		m.modal.SearchErr = msg.err.Error()

	case tea.KeyMsg:
		if m.modal.Type != ModalNone {
			switch msg.String() {
			case "esc":
				if m.currentURL != "" {
					m.modal.Close()
					return m, nil
				}
				return m, nil

			case "ctrl+c":
				return m, tea.Quit

			case "up", "ctrl+p":
				m.modal.Prev()
				return m, nil

			case "down", "ctrl+n":
				m.modal.Next()
				return m, nil

			case "tab":
				m.modal.Next()
				return m, nil

			case "shift+tab":
				m.modal.Prev()
				return m, nil

			case "enter":
				if m.modal.Type == ModalSearch {
					val := strings.TrimSpace(m.modal.Input.Value())
					if val == "" {
						return m, nil
					}

					if IsURL(val) {
						m.modal.Close()
						return m, m.navigate(val)
					}

					if len(m.modal.Filtered) > 0 && (val == m.modal.LastQuery || m.modal.Cursor > 0) {
						if sel := m.modal.SelectedItem(); sel != nil && sel.URL != "" {
							target := sel.URL
							m.modal.Close()
							return m, m.navigate(target)
						}
					}

					m.modal.IsSearching = true
					m.modal.SearchErr = ""
					return m, tea.Batch(m.spinner.Tick, m.searchQuery(val))
				}

				if m.modal.Type == ModalLinks {
					if sel := m.modal.SelectedItem(); sel != nil && sel.URL != "" {
						target := sel.URL
						m.modal.Close()
						return m, m.navigate(target)
					}
					return m, nil
				}

			default:
				var tiCmd tea.Cmd
				m.modal.Input, tiCmd = m.modal.Input.Update(msg)
				if m.modal.Type == ModalLinks {
					m.modal.Filter()
				}
				return m, tiCmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+o", "ctrl+g", "/", "o":
			m.modal.OpenSearch()
			return m, textinput.Blink

		case "f":
			if len(m.pageLinks) > 0 {
				m.modal.OpenLinks(m.pageLinks)
				return m, textinput.Blink
			}

		case "b":
			if len(m.history) > 1 {
				m.history = m.history[:len(m.history)-1]
				prevURL := m.history[len(m.history)-1]
				return m, m.fetchURL(prevURL)
			}

		case "r":
			if m.currentURL != "" {
				return m, m.fetchURL(m.currentURL)
			}

		case "g":
			m.viewport.GotoTop()

		case "G":
			m.viewport.GotoBottom()

		default:
			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *BrowserModel) navigate(targetURL string) tea.Cmd {
	if targetURL == "" {
		return nil
	}
	if len(m.history) == 0 || m.history[len(m.history)-1] != targetURL {
		m.history = append(m.history, targetURL)
	}
	return m.fetchURL(targetURL)
}

func (m *BrowserModel) fetchURL(targetURL string) tea.Cmd {
	m.isLoading = true
	m.loadingMsg = "Fetching " + targetURL + "..."
	return func() tea.Msg {
		start := time.Now()
		item, err := m.client.Fetch(context.Background(), targetURL)
		if err != nil {
			return fetchErrorMsg{err: err}
		}
		return fetchSuccessMsg{
			item:    item,
			latency: time.Since(start),
		}
	}
}

func (m *BrowserModel) searchQuery(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.client.Search(context.Background(), query)
		if err != nil {
			return searchErrorMsg{err: err}
		}
		return searchSuccessMsg{results: results}
	}
}

func (m *BrowserModel) renderMarkdown() {
	wrapWidth := m.width - 4
	if wrapWidth < 40 {
		wrapWidth = 40
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(OneDarkStyleConfig),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		m.renderedMD = m.rawMarkdown
		return
	}

	cleaned := NormalizeMarkdown(m.rawMarkdown)
	out, err := renderer.Render(cleaned)
	if err != nil {
		m.renderedMD = m.rawMarkdown
		return
	}

	m.renderedMD = out
}

func (m BrowserModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.modal.Type != ModalNone {
		return m.modal.RenderFullPage(m.width, m.height)
	}

	headerTag := HeaderTagStyle.Render("BROWSER")
	headerTitle := m.pageTitle
	if headerTitle == "" {
		headerTitle = "Web Viewer"
	}
	if len(headerTitle) > 40 {
		headerTitle = headerTitle[:37] + "..."
	}
	styledTitle := HeaderTitleStyle.Render(" " + headerTitle)
	leftHeader := lipgloss.JoinHorizontal(lipgloss.Left, headerTag, styledTitle)

	urlDisplay := m.currentURL
	if urlDisplay == "" {
		urlDisplay = "Press Ctrl+O to search or enter URL"
	}
	maxURLLen := m.width - lipgloss.Width(leftHeader) - 6
	if maxURLLen > 0 && len(urlDisplay) > maxURLLen {
		urlDisplay = urlDisplay[:maxURLLen-3] + "..."
	}
	rightHeader := HeaderURLStyle.Render(urlDisplay)

	spacing := m.width - lipgloss.Width(leftHeader) - lipgloss.Width(rightHeader) - 2
	if spacing < 1 {
		spacing = 1
	}
	headerBar := HeaderBarStyle.Width(m.width).Render(leftHeader + strings.Repeat(" ", spacing) + rightHeader)

	var contentArea string
	if m.isLoading {
		spinText := fmt.Sprintf(" %s %s", m.spinner.View(), m.loadingMsg)
		contentArea = lipgloss.Place(
			m.width,
			m.viewport.Height,
			lipgloss.Center,
			lipgloss.Center,
			SpinnerStyle.Render(spinText),
		)
	} else {
		contentArea = m.viewport.View()
	}

	statusLeft := lipgloss.JoinHorizontal(
		lipgloss.Left,
		StatusKeyStyle.Render(" Ctrl+O "), StatusValStyle.Render("Search  "),
		StatusKeyStyle.Render("f "), StatusValStyle.Render(fmt.Sprintf("Links (%d)  ", len(m.pageLinks))),
		StatusKeyStyle.Render("b "), StatusValStyle.Render("Back  "),
		StatusKeyStyle.Render("r "), StatusValStyle.Render("Reload  "),
		StatusKeyStyle.Render("j/k "), StatusValStyle.Render("Scroll  "),
		StatusKeyStyle.Render("q "), StatusValStyle.Render("Quit"),
	)

	scrollPercent := fmt.Sprintf("%3.0f%%", m.viewport.ScrollPercent()*100)
	statusRight := StatusValStyle.Render(scrollPercent + " ")
	statusSpacing := m.width - lipgloss.Width(statusLeft) - lipgloss.Width(statusRight) - 2
	if statusSpacing < 1 {
		statusSpacing = 1
	}
	statusBar := StatusBarStyle.Width(m.width).Render(statusLeft + strings.Repeat(" ", statusSpacing) + statusRight)

	return lipgloss.JoinVertical(lipgloss.Left, headerBar, contentArea, statusBar)
}

func IsURL(input string) bool {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return true
	}
	if !strings.Contains(input, " ") && strings.Contains(input, ".") {
		parts := strings.Split(input, "/")
		domain := parts[0]
		if strings.Contains(domain, ".") && len(domain) > 3 {
			if _, err := url.Parse("https://" + input); err == nil {
				return true
			}
		}
	}
	return false
}

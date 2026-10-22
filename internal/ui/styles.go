package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// One Dark Color Palette
	BgDark      = lipgloss.Color("#282C34")
	BgGutter    = lipgloss.Color("#21252B")
	BgSurface   = lipgloss.Color("#2C313A")
	BgSelection = lipgloss.Color("#3E4451")

	TextNormal = lipgloss.Color("#ABB2BF")
	TextWhite  = lipgloss.Color("#ECEFF4")
	TextMuted  = lipgloss.Color("#5C6370")

	TextMutedStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	AccentBlue   = lipgloss.Color("#61AFEF")
	AccentCyan   = lipgloss.Color("#56B6C2")
	AccentGreen  = lipgloss.Color("#98C379")
	AccentPurple = lipgloss.Color("#C678DD")
	AccentYellow = lipgloss.Color("#E5C07B")
	AccentRed    = lipgloss.Color("#E06C75")

	// Header Styles
	HeaderBarStyle = lipgloss.NewStyle().
			Background(BgGutter).
			Foreground(TextNormal).
			Padding(0, 1)

	HeaderTagStyle = lipgloss.NewStyle().
			Background(AccentPurple).
			Foreground(BgGutter).
			Bold(true).
			Padding(0, 1)

	HeaderTitleStyle = lipgloss.NewStyle().
				Foreground(TextWhite).
				Bold(true)

	HeaderURLStyle = lipgloss.NewStyle().
			Foreground(AccentBlue)

	// Status Bar Styles
	StatusBarStyle = lipgloss.NewStyle().
			Background(BgGutter).
			Foreground(TextMuted).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(AccentBlue).
			Bold(true)

	StatusValStyle = lipgloss.NewStyle().
			Foreground(TextNormal)

	// Search & Omnibox Bar
	SearchContainerStyle = lipgloss.NewStyle().
				Background(BgSurface).
				Padding(0, 1)

	SearchPromptStyle = lipgloss.NewStyle().
				Foreground(AccentGreen).
				Bold(true)

	CounterBadgeStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	// Modal Dialog Styles
	ModalBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentBlue).
			Background(BgGutter).
			Padding(0, 1)

	ModalTitleStyle = lipgloss.NewStyle().
			Foreground(TextWhite).
			Bold(true)

	ModalDividerStyle = lipgloss.NewStyle().
				Foreground(BgSelection)

	ModalFooterStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	// List & Results Styles
	SelectedRowStyle = lipgloss.NewStyle().
				Background(BgSelection).
				Foreground(TextWhite).
				Bold(true)

	NormalRowStyle = lipgloss.NewStyle().
			Foreground(TextWhite)

	ResultURLStyle = lipgloss.NewStyle().
			Foreground(AccentCyan)

	ResultSnippetStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(AccentBlue).
				Bold(true).
				Padding(0, 1)

	DividerStyle = lipgloss.NewStyle().
			Foreground(BgSurface)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(AccentRed).
			Bold(true)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(AccentBlue)
)

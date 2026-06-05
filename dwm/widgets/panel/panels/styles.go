package panels

import "github.com/charmbracelet/lipgloss"

const (
	base     = "#24273a"
	text     = "#cad3f5"
	subtext  = "#a5adcb"
	overlay0 = "#6e738d"
	mauve    = "#c6a0f6"
	peach    = "#f5a97f"
	green    = "#a6da95"
	red      = "#ed8796"
)

var (
	Page      = lipgloss.NewStyle().Background(lipgloss.Color(base)).Foreground(lipgloss.Color(text)).Padding(1, 2)
	title     = lipgloss.NewStyle().Foreground(lipgloss.Color(mauve)).Bold(true)
	muted     = lipgloss.NewStyle().Foreground(lipgloss.Color(subtext))
	outside   = lipgloss.NewStyle().Foreground(lipgloss.Color(overlay0))
	selected  = lipgloss.NewStyle().Foreground(lipgloss.Color(base)).Background(lipgloss.Color(mauve)).Bold(true)
	selectedP = lipgloss.NewStyle().Foreground(lipgloss.Color(base)).Background(lipgloss.Color(mauve)).Bold(true).Padding(0, 1)
	todaySt   = lipgloss.NewStyle().Foreground(lipgloss.Color(peach)).Bold(true).Underline(true)
	eventDay  = lipgloss.NewStyle().Foreground(lipgloss.Color(green)).Bold(true)
	ok        = lipgloss.NewStyle().Foreground(lipgloss.Color(green))
	bad       = lipgloss.NewStyle().Foreground(lipgloss.Color(red))
)

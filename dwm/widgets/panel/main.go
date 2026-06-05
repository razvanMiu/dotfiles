package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"dwm-widgets/panels"

	tea "github.com/charmbracelet/bubbletea"
)

type tab int

const (
	tabCalendar tab = iota
	tabAudio
)

type tickMsg struct{}
type switchTabMsg tab

type model struct {
	tab    tab
	width  int
	height int
	cal    panels.Calendar
	audio  panels.Audio
}

func main() {
	tuneRuntime()
	m := model{tab: readLastTab(), cal: panels.NewCalendar(), audio: panels.NewAudio()}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	go serveTabSwitches(p)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func tuneRuntime() {
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(50)
	debug.SetMemoryLimit(48 * 1024 * 1024)
}

func (m model) Init() tea.Cmd { return tick() }
func tick() tea.Cmd           { return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return tickMsg{} }) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case switchTabMsg:
		m.tab = tab(msg)
		writeLastTab(m.tab)
		if m.tab == tabAudio {
			m.audio.Refresh()
		} else {
			m.cal.Reload()
		}
		return m, nil
	case tickMsg:
		m.audio.Refresh()
		m.cal.Reload()
		return m, tick()
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.cal.SetSize(msg.Width, msg.Height)
		m.audio.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	if m.tab == tabCalendar {
		m.cal.Update(msg)
	} else {
		m.audio.Update(msg)
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	w := max(1, m.width-4)
	h := max(1, m.height-2)
	if m.tab == tabCalendar {
		return panels.Page.Width(w).Height(h).Render(m.cal.View(w, h))
	}
	return panels.Page.Width(w).Height(h).Render(m.audio.View(w, h))
}

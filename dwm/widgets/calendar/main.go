package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type event struct {
	date       string
	start      string
	end        string
	provider   string
	title      string
	body       string
	location   string
	meetingURL string
	allDay     bool
}

type eventJSON struct {
	Date       string `json:"date"`
	Start      string `json:"start"`
	End        string `json:"end"`
	AllDay     bool   `json:"all_day"`
	Provider   string `json:"provider"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Location   string `json:"location"`
	MeetingURL string `json:"meeting_url"`
}

type tickMsg struct{}

type model struct {
	selected  time.Time
	events    map[string][]event
	width     int
	height    int
	eventOff  int
	eventSel  int
	detailOff int
	cellW     int
}

const (
	calendarCols  = 7
	cellWidth     = 6
	calendarWidth = calendarCols * cellWidth
	contentWidth  = 44
	visibleEvents = 5
	detailLines   = 5

	base     = "#24273a"
	surface0 = "#363a4f"
	text     = "#cad3f5"
	subtext  = "#a5adcb"
	overlay0 = "#6e738d"
	mauve    = "#c6a0f6"
	peach    = "#f5a97f"
	green    = "#a6da95"
)

var (
	page      = lipgloss.NewStyle().Background(lipgloss.Color(base)).Foreground(lipgloss.Color(text)).Padding(1, 0)
	title     = lipgloss.NewStyle().Foreground(lipgloss.Color(mauve)).Bold(true)
	muted     = lipgloss.NewStyle().Foreground(lipgloss.Color(subtext))
	outside   = lipgloss.NewStyle().Foreground(lipgloss.Color(overlay0))
	todaySt   = lipgloss.NewStyle().Foreground(lipgloss.Color(peach)).Bold(true).Underline(true)
	selected  = lipgloss.NewStyle().Foreground(lipgloss.Color(base)).Background(lipgloss.Color(mauve)).Bold(true).Padding(0, 1)
	eventDay  = lipgloss.NewStyle().Foreground(lipgloss.Color(green)).Bold(true)
	eventHead = lipgloss.NewStyle().Foreground(lipgloss.Color(mauve)).Bold(true)
)

func main() {
	tuneRuntime()
	now := day(time.Now())
	m := model{selected: now, events: loadEvents(), eventSel: -1, cellW: cellWidth}
	m.autoSelect()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd { return reloadTick() }

func tuneRuntime() {
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(50)
	debug.SetMemoryLimit(32 * 1024 * 1024)
}

func reloadTick() tea.Cmd {
	return tea.Tick(10*time.Second, func(time.Time) tea.Msg { return tickMsg{} })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.events = loadEvents()
		m.clampSelection()
		return m, reloadTick()
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "h", "left":
			m.shiftDay(-1)
		case "l", "right":
			m.shiftDay(1)
		case "j", "down":
			m.shiftDay(7)
		case "k", "up":
			m.shiftDay(-7)
		case "H", "pgup":
			m.shiftMonth(-1)
		case "L", "pgdown":
			m.shiftMonth(1)
		case "[", "ctrl+u":
			m.scrollDetails(-1)
		case "]", "ctrl+d":
			m.scrollDetails(1)
		case "t":
			m.selected = day(time.Now())
			m.eventOff = 0
			m.autoSelect()
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp {
			if _, ok := m.dateAt(msg.X, msg.Y); ok {
				m.shiftMonth(-1)
			} else if m.inDetails(msg.Y) {
				m.scrollDetails(-1)
			} else {
				m.scrollEvents(-1)
			}
			break
		}
		if msg.Button == tea.MouseButtonWheelDown {
			if _, ok := m.dateAt(msg.X, msg.Y); ok {
				m.shiftMonth(1)
			} else if m.inDetails(msg.Y) {
				m.scrollDetails(1)
			} else {
				m.scrollEvents(1)
			}
			break
		}
		if msg.Button == tea.MouseButtonLeft && msg.Action != tea.MouseActionMotion {
			if d, ok := m.dateAt(msg.X, msg.Y); ok {
				m.selected = d
				m.eventOff = 0
				m.autoSelect()
			} else if idx := m.eventAt(msg.X, msg.Y); idx >= 0 {
				m.eventSel = idx
				m.detailOff = 0
			} else if msg.Y == eventHeaderY()+1 && m.eventOff > 0 {
				m.scrollEvents(-1)
			} else if msg.Y == eventRowsY()+visibleEvents {
				m.scrollEvents(1)
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}
	innerW := max(calendarWidth, m.width)
	availableLines := max(1, m.height-2)

	var body strings.Builder
	body.WriteString(center(title.Render(m.selected.Format("January 2006")), innerW))
	body.WriteString("\n\n")
	body.WriteString(m.renderCalendar(innerW))
	body.WriteString("\n")
	body.WriteString(m.renderEvents(innerW))

	footer := center(muted.Render("h/l day · j/k week · H/L mon · [/] det · t/q"), innerW)
	bodyLines := strings.Count(body.String(), "\n") + 1
	blankLines := max(1, availableLines-bodyLines-1)

	return page.Width(innerW).Render(body.String() + strings.Repeat("\n", blankLines) + footer)
}

func (m model) renderCalendar(w int) string {
	var b strings.Builder
	days := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	var row strings.Builder
	for _, d := range days {
		row.WriteString(muted.Width(m.cellW).Align(lipgloss.Center).Render(d))
	}
	b.WriteString(center(row.String(), w))
	b.WriteByte('\n')
	weeks := monthWeeks(m.selected)
	today := day(time.Now())
	for _, week := range weeks {
		row.Reset()
		for _, d := range week {
			label := fmt.Sprintf("%02d", d.Day())
			st := lipgloss.NewStyle().Width(m.cellW).Align(lipgloss.Center).Foreground(lipgloss.Color(text))
			if d.Month() != m.selected.Month() {
				st = outside.Width(m.cellW).Align(lipgloss.Center)
			}
			if len(m.events[dateKey(d)]) > 0 && d.Month() == m.selected.Month() {
				label = eventDay.Render(label)
			}
			if sameDay(d, today) {
				label = todaySt.Render(fmt.Sprintf("%02d", d.Day()))
			}
			if sameDay(d, m.selected) {
				label = selected.Render(fmt.Sprintf("%02d", d.Day()))
			}
			row.WriteString(st.Render(label))
		}
		b.WriteString(center(row.String(), w))
		b.WriteByte('\n')
	}
	return b.String()
}

func (m model) renderEvents(w int) string {
	var b strings.Builder
	blockW := min(contentWidth, w)
	indent := strings.Repeat(" ", max(0, (w-blockW)/2))
	lineW := blockW
	b.WriteString(indent + eventHead.Render("Events — "+m.selected.Format("Mon 02 Jan")) + "\n")
	evs := m.dayEvents()
	if len(evs) == 0 {
		b.WriteString(indent + muted.Render("No events") + strings.Repeat("\n", visibleEvents+1))
		b.WriteString("\n" + indent + eventHead.Render("Details") + "\n")
		b.WriteString(indent + muted.Render("No event selected") + strings.Repeat("\n", detailLines-1))
		return b.String()
	}
	maxOff := max(0, len(evs)-visibleEvents)
	if m.eventOff > maxOff {
		m.eventOff = maxOff
	}
	if m.eventOff > 0 {
		b.WriteString(indent + muted.Render("↑ more") + "\n")
	} else {
		b.WriteString("\n")
	}
	shown := evs[m.eventOff:min(len(evs), m.eventOff+visibleEvents)]
	for i := 0; i < visibleEvents; i++ {
		if i >= len(shown) {
			b.WriteString("\n")
			continue
		}
		idx := m.eventOff + i
		line := ellipsize(eventLine(shown[i]), lineW)
		if m.eventSel == idx {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color(base)).Background(lipgloss.Color(mauve)).Render(line)
		}
		b.WriteString(indent + line + "\n")
	}
	if m.eventOff+visibleEvents < len(evs) {
		b.WriteString(indent + muted.Render("↓ more") + "\n")
	} else {
		b.WriteString("\n")
	}
	b.WriteString("\n" + indent + eventHead.Render("Details") + "\n")
	b.WriteString(m.renderDetails(indent, lineW, evs))
	return b.String()
}

func (m model) renderDetails(indent string, lineW int, evs []event) string {
	if m.eventSel < 0 || m.eventSel >= len(evs) {
		return indent + muted.Render("Select an event") + strings.Repeat("\n", detailLines-1)
	}
	e := evs[m.eventSel]
	var lines []string
	lines = append(lines, wrapLine(fmt.Sprintf("%s · %s", e.title, e.provider), lineW)...)
	for _, line := range detailText(e) {
		lines = append(lines, wrapLine(line, lineW)...)
	}
	maxOff := max(0, len(lines)-detailLines)
	if m.detailOff > maxOff {
		m.detailOff = maxOff
	}
	var b strings.Builder
	for i := 0; i < detailLines; i++ {
		idx := m.detailOff + i
		if idx < len(lines) {
			b.WriteString(indent + muted.Render(ellipsize(lines[idx], lineW)))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (m *model) shiftDay(n int) {
	m.selected = day(m.selected.AddDate(0, 0, n))
	m.eventOff = 0
	m.autoSelect()
}
func (m *model) shiftMonth(n int) {
	m.selected = day(m.selected.AddDate(0, n, 0))
	m.eventOff = 0
	m.autoSelect()
}
func (m *model) scrollEvents(n int) {
	maxOff := max(0, len(m.dayEvents())-visibleEvents)
	m.eventOff = max(0, min(maxOff, m.eventOff+n))
}
func (m *model) scrollDetails(n int) { m.detailOff = max(0, m.detailOff+n) }

func (m *model) autoSelect() {
	evs := m.dayEvents()
	m.detailOff = 0
	if len(evs) == 0 {
		m.eventSel = -1
		return
	}
	m.eventSel = 0
	if sameDay(m.selected, day(time.Now())) {
		now := time.Now().Format("15:04")
		for i, e := range evs {
			if !e.allDay && e.start >= now {
				m.eventSel = i
				break
			}
		}
	}
	if m.eventSel >= m.eventOff+visibleEvents {
		m.eventOff = max(0, m.eventSel-visibleEvents+1)
	}
}
func (m *model) clampSelection() {
	evs := m.dayEvents()
	if len(evs) == 0 {
		m.eventSel = -1
		return
	}
	if m.eventSel >= len(evs) {
		m.eventSel = len(evs) - 1
	}
}

func (m model) dayEvents() []event {
	evs := append([]event(nil), m.events[dateKey(m.selected)]...)
	sort.Slice(evs, func(i, j int) bool {
		if evs[i].allDay != evs[j].allDay {
			return evs[i].allDay
		}
		return evs[i].start < evs[j].start
	})
	return evs
}

func eventHeaderY() int              { return 11 }
func eventRowsY() int                { return eventHeaderY() + 2 }
func detailsStartY() int             { return eventRowsY() + visibleEvents + 3 }
func (m model) inDetails(y int) bool { return y >= detailsStartY()+1 }
func (m model) eventAt(x, y int) int {
	left, right := m.contentBounds()
	if x < left || x >= right || y < eventRowsY() || y >= eventRowsY()+visibleEvents {
		return -1
	}
	idx := m.eventOff + y - eventRowsY()
	if idx >= len(m.dayEvents()) {
		return -1
	}
	return idx
}

func (m model) contentBounds() (int, int) {
	innerW := max(calendarWidth, m.width)
	blockW := min(contentWidth, innerW)
	left := max(0, (innerW-blockW)/2)
	return left, left + blockW
}

func eventLine(e event) string {
	when := e.start
	if e.allDay {
		when = "All"
	}
	return fmt.Sprintf("%-5s %s", when, e.title)
}
func detailText(e event) []string {
	var lines []string
	if strings.TrimSpace(e.body) != "" {
		lines = append(lines, strings.Split(strings.TrimSpace(e.body), "\n")...)
	}
	if e.location != "" {
		lines = append(lines, "Location: "+e.location)
	}
	if e.meetingURL != "" {
		lines = append(lines, "Meeting: "+e.meetingURL)
	}
	if len(lines) == 0 {
		lines = append(lines, "No details")
	}
	return lines
}
func wrapLine(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	line := ""
	for _, word := range words {
		if line == "" {
			line = word
			continue
		}
		if len([]rune(line))+1+len([]rune(word)) <= width {
			line += " " + word
		} else {
			lines = append(lines, line)
			line = word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func ellipsize(s string, width int) string {
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width <= 1 {
		return "…"
	}
	return string(r[:width-1]) + "…"
}

func (m model) dateAt(x, y int) (time.Time, bool) {
	innerW := max(calendarWidth, m.width)
	gridX := max(0, (innerW-calendarWidth)/2)
	gridY := 3
	row, col := y-gridY-1, (x-gridX)/m.cellW
	if row < 0 || row >= 6 || col < 0 || col >= 7 || x < gridX {
		return time.Time{}, false
	}
	weeks := monthWeeks(m.selected)
	if row >= len(weeks) {
		return time.Time{}, false
	}
	return weeks[row][col], true
}

func loadEvents() map[string][]event {
	out := map[string][]event{}
	cache := filepath.Join(os.Getenv("XDG_CACHE_HOME"), "dwm", "calendar-agenda")
	if os.Getenv("XDG_CACHE_HOME") == "" {
		cache = filepath.Join(os.Getenv("HOME"), ".cache", "dwm", "calendar-agenda")
	}
	data, err := os.ReadFile(cache)
	if err != nil {
		return out
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var j eventJSON
		if err := json.Unmarshal([]byte(line), &j); err != nil {
			continue
		}
		if _, err := time.Parse("2006-01-02", j.Date); err != nil {
			continue
		}
		e := event{date: j.Date, start: j.Start, end: j.End, provider: j.Provider, title: j.Title, body: j.Body, location: j.Location, meetingURL: j.MeetingURL, allDay: j.AllDay}
		out[j.Date] = append(out[j.Date], e)
	}
	return out
}

func monthWeeks(t time.Time) [][]time.Time {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	start := first.AddDate(0, 0, -(int(first.Weekday())+6)%7)
	weeks := [][]time.Time{}
	for r := 0; r < 6; r++ {
		week := []time.Time{}
		for c := 0; c < 7; c++ {
			week = append(week, start.AddDate(0, 0, r*7+c))
		}
		weeks = append(weeks, week)
	}
	return weeks
}

func day(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
func sameDay(a, b time.Time) bool { return a.Year() == b.Year() && a.YearDay() == b.YearDay() }
func dateKey(t time.Time) string  { return t.Format("2006-01-02") }
func center(s string, w int) string {
	return lipgloss.NewStyle().Width(w).Align(lipgloss.Center).Render(s)
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

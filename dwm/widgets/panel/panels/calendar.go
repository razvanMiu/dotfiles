package panels

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type event struct {
	date, start, end, provider, title, body, location, meetingURL string
	allDay                                                        bool
}
type eventJSON struct {
	Date, Start, End, Provider, Title, Body, Location string
	AllDay                                            bool   `json:"all_day"`
	MeetingURL                                        string `json:"meeting_url"`
}
type Calendar struct {
	selected           time.Time
	events             map[string][]event
	width, height      int
	eventOff, eventSel int
	detailOff          int
}

const (
	calendarCols  = 7
	cellWidth     = 6
	calendarWidth = calendarCols * cellWidth
	contentWidth  = 44
	visibleEvents = 5
	detailLines   = 5
)

func NewCalendar() Calendar {
	m := Calendar{selected: day(time.Now()), events: loadEvents(), eventSel: -1}
	m.autoSelect()
	return m
}
func (m Calendar) WidthHint() int             { return contentWidth }
func (m Calendar) HeightHint() int            { return 30 }
func (m *Calendar) SetSize(width, height int) { m.width, m.height = width, height }
func (m *Calendar) Reload()                   { m.events = loadEvents(); m.clampSelection() }
func (m *Calendar) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
		}
		if msg.Button == tea.MouseButtonWheelDown {
			if _, ok := m.dateAt(msg.X, msg.Y); ok {
				m.shiftMonth(1)
			} else if m.inDetails(msg.Y) {
				m.scrollDetails(1)
			} else {
				m.scrollEvents(1)
			}
		}
		if msg.Button == tea.MouseButtonLeft && msg.Action != tea.MouseActionMotion {
			if d, ok := m.dateAt(msg.X, msg.Y); ok {
				m.selected = d
				m.eventOff = 0
				m.autoSelect()
			} else if idx := m.eventAt(msg.X, msg.Y); idx >= 0 {
				m.eventSel = idx
				m.detailOff = 0
			}
		}
	}
}
func (m Calendar) View(w, h int) string {
	innerW := max(1, w)
	var body strings.Builder
	body.WriteString(center(title.Render(m.selected.Format("January 2006")), innerW))
	body.WriteString("\n\n")
	body.WriteString(m.renderCalendar(innerW))
	body.WriteByte('\n')
	body.WriteString(m.renderEvents(innerW))
	lines := strings.Split(strings.TrimRight(body.String(), "\n"), "\n")
	footer := center(muted.Render("h/l day · j/k week · H/L mon · [/] det · t/q"), innerW)
	if h <= 1 {
		return fitLines(lines, w, h)
	}
	lines = lines[:min(len(lines), h-1)]
	for len(lines) < h-1 {
		lines = append(lines, "")
	}
	lines = append(lines, footer)
	return fitLines(lines, w, h)
}
func (m Calendar) renderCalendar(w int) string {
	var b, row strings.Builder
	for _, d := range []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"} {
		row.WriteString(muted.Width(cellWidth).Align(lipgloss.Center).Render(d))
	}
	b.WriteString(center(row.String(), w))
	b.WriteByte('\n')
	today := day(time.Now())
	for _, week := range monthWeeks(m.selected) {
		row.Reset()
		for _, d := range week {
			label := fmt.Sprintf("%02d", d.Day())
			st := lipgloss.NewStyle().Width(cellWidth).Align(lipgloss.Center).Foreground(lipgloss.Color(text))
			if d.Month() != m.selected.Month() {
				st = outside.Width(cellWidth).Align(lipgloss.Center)
			}
			if len(m.events[dateKey(d)]) > 0 && d.Month() == m.selected.Month() {
				label = eventDay.Render(label)
			}
			if sameDay(d, today) {
				label = todaySt.Render(fmt.Sprintf("%02d", d.Day()))
			}
			if sameDay(d, m.selected) {
				label = selectedP.Render(fmt.Sprintf("%02d", d.Day()))
			}
			row.WriteString(st.Render(label))
		}
		b.WriteString(center(row.String(), w))
		b.WriteByte('\n')
	}
	return b.String()
}
func (m Calendar) renderEvents(w int) string {
	var b strings.Builder
	blockW := min(contentWidth, w)
	indent := strings.Repeat(" ", max(0, (w-blockW)/2))
	lineW := blockW
	b.WriteString(indent + title.Render("Events — "+m.selected.Format("Mon 02 Jan")) + "\n")
	evs := m.dayEvents()
	if len(evs) == 0 {
		return b.String() + indent + muted.Render("No events") + strings.Repeat("\n", visibleEvents+1) + "\n" + indent + title.Render("Details") + "\n" + indent + muted.Render("No event selected") + strings.Repeat("\n", detailLines-1)
	}
	if m.eventOff > max(0, len(evs)-visibleEvents) {
		m.eventOff = max(0, len(evs)-visibleEvents)
	}
	if m.eventOff > 0 {
		b.WriteString(indent + muted.Render("↑ more") + "\n")
	} else {
		b.WriteByte('\n')
	}
	shown := evs[m.eventOff:min(len(evs), m.eventOff+visibleEvents)]
	for i := 0; i < visibleEvents; i++ {
		if i >= len(shown) {
			b.WriteByte('\n')
			continue
		}
		idx := m.eventOff + i
		line := ellipsize(eventLine(shown[i]), lineW)
		if m.eventSel == idx {
			line = selected.Render(line)
		}
		b.WriteString(indent + line + "\n")
	}
	if m.eventOff+visibleEvents < len(evs) {
		b.WriteString(indent + muted.Render("↓ more") + "\n")
	} else {
		b.WriteByte('\n')
	}
	b.WriteString("\n" + indent + title.Render("Details") + "\n")
	b.WriteString(m.renderDetails(indent, lineW, evs))
	return b.String()
}
func (m Calendar) renderDetails(indent string, lineW int, evs []event) string {
	if m.eventSel < 0 || m.eventSel >= len(evs) {
		return indent + muted.Render("Select an event") + strings.Repeat("\n", detailLines-1)
	}
	e := evs[m.eventSel]
	lines := wrapLine(fmt.Sprintf("%s · %s", e.title, e.provider), lineW)
	for _, line := range detailText(e) {
		lines = append(lines, wrapLine(line, lineW)...)
	}
	if m.detailOff > max(0, len(lines)-detailLines) {
		m.detailOff = max(0, len(lines)-detailLines)
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
func (m *Calendar) shiftDay(n int) {
	m.selected = day(m.selected.AddDate(0, 0, n))
	m.eventOff = 0
	m.autoSelect()
}
func (m *Calendar) shiftMonth(n int) {
	m.selected = day(m.selected.AddDate(0, n, 0))
	m.eventOff = 0
	m.autoSelect()
}
func (m *Calendar) scrollEvents(n int) {
	m.eventOff = max(0, min(max(0, len(m.dayEvents())-visibleEvents), m.eventOff+n))
}
func (m *Calendar) scrollDetails(n int) { m.detailOff = max(0, m.detailOff+n) }
func (m *Calendar) autoSelect() {
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
func (m *Calendar) clampSelection() {
	evs := m.dayEvents()
	if len(evs) == 0 {
		m.eventSel = -1
		return
	}
	if m.eventSel >= len(evs) {
		m.eventSel = len(evs) - 1
	}
}
func (m Calendar) dayEvents() []event {
	evs := append([]event(nil), m.events[dateKey(m.selected)]...)
	sort.Slice(evs, func(i, j int) bool {
		if evs[i].allDay != evs[j].allDay {
			return evs[i].allDay
		}
		return evs[i].start < evs[j].start
	})
	return evs
}
func eventHeaderY() int                 { return 13 }
func eventRowsY() int                   { return eventHeaderY() + 2 }
func detailsStartY() int                { return eventRowsY() + visibleEvents + 3 }
func (m Calendar) inDetails(y int) bool { return y >= detailsStartY()+1 }
func (m Calendar) eventAt(x, y int) int {
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
func (m Calendar) contentBounds() (int, int) {
	innerW := max(calendarWidth, m.width)
	blockW := min(contentWidth, innerW)
	left := max(0, (innerW-blockW)/2)
	return left, left + blockW
}
func (m Calendar) dateAt(x, y int) (time.Time, bool) {
	innerW := max(calendarWidth, m.width)
	gridX := max(0, (innerW-calendarWidth)/2)
	gridY := 5
	row, col := y-gridY-1, (x-gridX)/cellWidth
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
		if json.Unmarshal([]byte(line), &j) != nil {
			continue
		}
		if _, err := time.Parse("2006-01-02", j.Date); err != nil {
			continue
		}
		out[j.Date] = append(out[j.Date], event{date: j.Date, start: j.Start, end: j.End, provider: j.Provider, title: j.Title, body: j.Body, location: j.Location, meetingURL: j.MeetingURL, allDay: j.AllDay})
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

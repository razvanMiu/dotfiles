package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	tabCalendar tab = iota
	tabAudio
)

type tickMsg struct{}
type reloadMsg struct{}
type switchTabMsg tab

type model struct {
	tab    tab
	width  int
	height int
	cal    calendarModel
	audio  audioModel
}

const (
	base     = "#24273a"
	surface0 = "#363a4f"
	text     = "#cad3f5"
	subtext  = "#a5adcb"
	overlay0 = "#6e738d"
	mauve    = "#c6a0f6"
	peach    = "#f5a97f"
	green    = "#a6da95"
	red      = "#ed8796"
	yellow   = "#eed49f"
)

var (
	page      = lipgloss.NewStyle().Background(lipgloss.Color(base)).Foreground(lipgloss.Color(text)).Padding(1, 2)
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

func main() {
	tuneRuntime()
	m := model{tab: readLastTab(), cal: newCalendar(), audio: newAudio()}
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
		if m.tab == tabAudio {
			m.audio.refresh()
		} else {
			m.cal.reload()
		}
		return m, nil
	case tickMsg:
		m.audio.refresh()
		m.cal.reload()
		return m, tick()
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.cal.width, m.cal.height = msg.Width, msg.Height
		m.audio.width, m.audio.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.tab = (m.tab + 1) % 2
			writeLastTab(m.tab)
			return m, nil
		case "shift+tab":
			m.tab = (m.tab + 1) % 2
			writeLastTab(m.tab)
			return m, nil
		}
	}
	if m.tab == tabCalendar {
		m.cal.update(msg)
	} else {
		m.audio.update(msg)
	}
	return m, nil
}

func (m model) View() string {
	w := max(58, m.width-4)
	var b strings.Builder
	b.WriteString(center(tabHeader(m.tab), w))
	b.WriteByte('\n')
	b.WriteByte('\n')
	if m.tab == tabCalendar {
		b.WriteString(m.cal.view(w, max(1, m.height-4)))
	} else {
		b.WriteString(m.audio.view(w))
	}
	return page.Width(w).Render(b.String())
}

func tabHeader(t tab) string {
	cal := title.Render("Calendar")
	aud := title.Render("Audio")
	if t == tabCalendar {
		cal = selectedP.Render("Calendar")
	} else {
		aud = selectedP.Render("Audio")
	}
	return cal + muted.Render("  Tab switch  ") + aud
}

func serveTabSwitches(p *tea.Program) {
	sock := socketPath()
	_ = os.Remove(sock)
	_ = os.MkdirAll(filepath.Dir(sock), 0755)
	addr, err := net.ResolveUnixAddr("unixgram", sock)
	if err != nil {
		return
	}
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		return
	}
	defer conn.Close()
	defer os.Remove(sock)
	buf := make([]byte, 64)
	for {
		n, _, err := conn.ReadFromUnix(buf)
		if err != nil {
			return
		}
		switch strings.TrimSpace(string(buf[:n])) {
		case "calendar":
			p.Send(switchTabMsg(tabCalendar))
		case "audio":
			p.Send(switchTabMsg(tabAudio))
		}
	}
}

func readLastTab() tab {
	switch strings.TrimSpace(readFile(lastTabPath())) {
	case "audio":
		return tabAudio
	default:
		return tabCalendar
	}
}
func writeLastTab(t tab) {
	v := "calendar"
	if t == tabAudio {
		v = "audio"
	}
	_ = os.MkdirAll(filepath.Dir(lastTabPath()), 0755)
	_ = os.WriteFile(lastTabPath(), []byte(v), 0644)
}
func readFile(path string) string { b, _ := os.ReadFile(path); return string(b) }
func stateDir() string {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return filepath.Join(v, "dwm")
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "state", "dwm")
}
func socketPath() string  { return filepath.Join(stateDir(), "widgets.sock") }
func lastTabPath() string { return filepath.Join(stateDir(), "last-widget-tab") }

/* Calendar */

type event struct {
	date, start, end, provider, title, body, location, meetingURL string
	allDay                                                        bool
}
type eventJSON struct {
	Date, Start, End, Provider, Title, Body, Location string
	AllDay                                            bool   `json:"all_day"`
	MeetingURL                                        string `json:"meeting_url"`
}
type calendarModel struct {
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

func newCalendar() calendarModel {
	m := calendarModel{selected: day(time.Now()), events: loadEvents(), eventSel: -1}
	m.autoSelect()
	return m
}
func (m *calendarModel) reload() { m.events = loadEvents(); m.clampSelection() }
func (m *calendarModel) update(msg tea.Msg) {
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
func (m calendarModel) view(w, h int) string {
	innerW := max(calendarWidth, w)
	var body strings.Builder
	body.WriteString(center(title.Render(m.selected.Format("January 2006")), innerW))
	body.WriteString("\n\n")
	body.WriteString(m.renderCalendar(innerW))
	body.WriteByte('\n')
	body.WriteString(m.renderEvents(innerW))
	footer := center(muted.Render("h/l day · j/k week · H/L mon · [/] det · t/q"), innerW)
	bodyLines := strings.Count(body.String(), "\n") + 1
	blankLines := max(1, h-bodyLines-1)
	return body.String() + strings.Repeat("\n", blankLines) + footer
}
func (m calendarModel) renderCalendar(w int) string {
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
func (m calendarModel) renderEvents(w int) string {
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
func (m calendarModel) renderDetails(indent string, lineW int, evs []event) string {
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
func (m *calendarModel) shiftDay(n int) {
	m.selected = day(m.selected.AddDate(0, 0, n))
	m.eventOff = 0
	m.autoSelect()
}
func (m *calendarModel) shiftMonth(n int) {
	m.selected = day(m.selected.AddDate(0, n, 0))
	m.eventOff = 0
	m.autoSelect()
}
func (m *calendarModel) scrollEvents(n int) {
	m.eventOff = max(0, min(max(0, len(m.dayEvents())-visibleEvents), m.eventOff+n))
}
func (m *calendarModel) scrollDetails(n int) { m.detailOff = max(0, m.detailOff+n) }
func (m *calendarModel) autoSelect() {
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
func (m *calendarModel) clampSelection() {
	evs := m.dayEvents()
	if len(evs) == 0 {
		m.eventSel = -1
		return
	}
	if m.eventSel >= len(evs) {
		m.eventSel = len(evs) - 1
	}
}
func (m calendarModel) dayEvents() []event {
	evs := append([]event(nil), m.events[dateKey(m.selected)]...)
	sort.Slice(evs, func(i, j int) bool {
		if evs[i].allDay != evs[j].allDay {
			return evs[i].allDay
		}
		return evs[i].start < evs[j].start
	})
	return evs
}
func eventHeaderY() int                      { return 13 }
func eventRowsY() int                        { return eventHeaderY() + 2 }
func detailsStartY() int                     { return eventRowsY() + visibleEvents + 3 }
func (m calendarModel) inDetails(y int) bool { return y >= detailsStartY()+1 }
func (m calendarModel) eventAt(x, y int) int {
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
func (m calendarModel) contentBounds() (int, int) {
	innerW := max(calendarWidth, m.width)
	blockW := min(contentWidth, innerW)
	left := max(0, (innerW-blockW)/2)
	return left, left + blockW
}
func (m calendarModel) dateAt(x, y int) (time.Time, bool) {
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

/* Audio */

type volumeMap map[string]struct {
	ValuePercent string `json:"value_percent"`
}
type device struct {
	Index             int `json:"index"`
	Name, Description string
	Mute              bool
	Volume            volumeMap
}
type stream struct {
	Index      int `json:"index"`
	Mute       bool
	Volume     volumeMap
	Properties map[string]string `json:"properties"`
}
type info struct {
	DefaultSinkName   string `json:"default_sink_name"`
	DefaultSourceName string `json:"default_source_name"`
}
type audioState struct {
	Info                  info
	Sinks, Sources        []device
	Inputs, SourceOutputs []stream
}
type audioModel struct {
	st              audioState
	section, cursor int
	width, height   int
	err             string
}

func newAudio() audioModel { m := audioModel{}; m.refresh(); return m }
func (m *audioModel) update(msg tea.Msg) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "tab":
			m.section = (m.section + 1) % 4
			m.cursor = 0
		case "shift+tab":
			m.section = (m.section + 3) % 4
			m.cursor = 0
		case "j", "down":
			m.cursor = min(m.cursor+1, m.sectionLen()-1)
		case "k", "up":
			m.cursor = max(0, m.cursor-1)
		case "h", "left":
			m.adjust(-5)
		case "l", "right":
			m.adjust(5)
		case "m":
			m.toggleMute()
		case "enter":
			m.activate()
		case "r":
			m.refresh()
		}
	}
}
func (m audioModel) view(w int) string {
	var b strings.Builder
	b.WriteString(center(title.Render("Audio"), w))
	b.WriteString("\n")
	b.WriteString(center(muted.Render("Tab section · j/k row · h/l vol · m mute · Enter default · r/q"), w))
	b.WriteString("\n\n")
	b.WriteString(m.renderSection(0, "Outputs", m.renderDevices(m.st.Sinks, m.st.Info.DefaultSinkName, true)))
	b.WriteString(m.renderSection(1, "Inputs", m.renderDevices(m.st.Sources, m.st.Info.DefaultSourceName, false)))
	b.WriteString(m.renderSection(2, "Apps", m.renderStreams(m.st.Inputs)))
	b.WriteString(m.renderSection(3, "Recording", m.renderStreams(m.st.SourceOutputs)))
	if m.err != "" {
		b.WriteString("\n" + bad.Render(m.err))
	}
	return b.String()
}
func (m audioModel) renderSection(i int, name string, lines []string) string {
	head := "  " + name
	if m.section == i {
		head = "▶ " + name
	}
	var b strings.Builder
	b.WriteString(title.Render(head))
	b.WriteString(muted.Render(fmt.Sprintf("  %d", len(lines))))
	b.WriteByte('\n')
	if len(lines) == 0 {
		b.WriteString("  " + muted.Render("No active streams") + "\n\n")
		return b.String()
	}
	for n, l := range lines {
		if m.section == i && m.cursor == n {
			l = selected.Render(l)
		}
		b.WriteString("  " + l + "\n")
	}
	b.WriteByte('\n')
	return b.String()
}
func (m audioModel) renderDevices(ds []device, def string, output bool) []string {
	var out []string
	for _, d := range ds {
		if !output && strings.Contains(d.Name, ".monitor") {
			continue
		}
		mark := " "
		if d.Name == def {
			mark = ok.Render("●")
		}
		mute := ""
		if d.Mute {
			mute = " " + bad.Render("muted")
		}
		p := percent(d.Volume)
		out = append(out, fmt.Sprintf("%s %-30s %3d%% %s%s", mark, short(d.Description, 30), p, bar(p, 16), mute))
	}
	return out
}
func (m audioModel) renderStreams(ss []stream) []string {
	var out []string
	sort.Slice(ss, func(i, j int) bool { return appName(ss[i]) < appName(ss[j]) })
	for _, s := range ss {
		mute := ""
		if s.Mute {
			mute = " " + bad.Render("muted")
		}
		p := percent(s.Volume)
		out = append(out, fmt.Sprintf("%-32s %3d%% %s%s", short(appName(s), 32), p, bar(p, 16), mute))
	}
	return out
}
func (m *audioModel) refresh() {
	st, err := loadAudio()
	if err != nil {
		m.err = err.Error()
		return
	}
	m.st = st
	m.err = ""
	if m.cursor >= m.sectionLen() {
		m.cursor = max(0, m.sectionLen()-1)
	}
}
func loadAudio() (audioState, error) {
	var st audioState
	if err := jsonCmd(&st.Info, "pactl", "--format=json", "info"); err != nil {
		return st, err
	}
	_ = jsonCmd(&st.Sinks, "pactl", "--format=json", "list", "sinks")
	_ = jsonCmd(&st.Sources, "pactl", "--format=json", "list", "sources")
	_ = jsonCmd(&st.Inputs, "pactl", "--format=json", "list", "sink-inputs")
	_ = jsonCmd(&st.SourceOutputs, "pactl", "--format=json", "list", "source-outputs")
	return st, nil
}
func jsonCmd(v any, name string, args ...string) error {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return err
	}
	return json.Unmarshal(out, v)
}
func (m audioModel) sectionLen() int {
	switch m.section {
	case 0:
		return len(m.renderDevices(m.st.Sinks, m.st.Info.DefaultSinkName, true))
	case 1:
		return len(filterSources(m.st.Sources))
	case 2:
		return len(m.st.Inputs)
	case 3:
		return len(m.st.SourceOutputs)
	}
	return 0
}
func (m audioModel) selectedID() string {
	switch m.section {
	case 0:
		if m.cursor < len(m.st.Sinks) {
			return m.st.Sinks[m.cursor].Name
		}
	case 1:
		ds := filterSources(m.st.Sources)
		if m.cursor < len(ds) {
			return ds[m.cursor].Name
		}
	case 2:
		if m.cursor < len(m.st.Inputs) {
			return fmt.Sprint(m.st.Inputs[m.cursor].Index)
		}
	case 3:
		if m.cursor < len(m.st.SourceOutputs) {
			return fmt.Sprint(m.st.SourceOutputs[m.cursor].Index)
		}
	}
	return ""
}
func (m *audioModel) adjust(delta int) {
	id := m.selectedID()
	if id == "" {
		return
	}
	switch m.section {
	case 0:
		exec.Command("pactl", "set-sink-volume", id, fmt.Sprintf("%+d%%", delta)).Run()
	case 1:
		exec.Command("pactl", "set-source-volume", id, fmt.Sprintf("%+d%%", delta)).Run()
	case 2:
		exec.Command("pactl", "set-sink-input-volume", id, fmt.Sprintf("%+d%%", delta)).Run()
	case 3:
		exec.Command("pactl", "set-source-output-volume", id, fmt.Sprintf("%+d%%", delta)).Run()
	}
	m.refresh()
	refreshStatusBar()
}
func (m *audioModel) toggleMute() {
	id := m.selectedID()
	if id == "" {
		return
	}
	cmds := []string{"set-sink-mute", "set-source-mute", "set-sink-input-mute", "set-source-output-mute"}
	exec.Command("pactl", cmds[m.section], id, "toggle").Run()
	m.refresh()
	refreshStatusBar()
}
func (m *audioModel) activate() {
	id := m.selectedID()
	if id == "" {
		return
	}
	if m.section == 0 {
		exec.Command("pactl", "set-default-sink", id).Run()
	}
	if m.section == 1 {
		exec.Command("pactl", "set-default-source", id).Run()
	}
	m.refresh()
	refreshStatusBar()
}
func refreshStatusBar() {
	if cmd := dwmStatusPath(); cmd != "" {
		exec.Command(cmd, "--refresh", "volume").Run()
	}
}
func dwmStatusPath() string {
	if home := os.Getenv("HOME"); home != "" {
		path := home + "/.local/bin/dwm-status"
		if st, err := os.Stat(path); err == nil && !st.IsDir() && st.Mode()&0111 != 0 {
			return path
		}
	}
	if path, err := exec.LookPath("dwm-status"); err == nil {
		return path
	}
	return ""
}
func filterSources(in []device) []device {
	var out []device
	for _, d := range in {
		if !strings.Contains(d.Name, ".monitor") {
			out = append(out, d)
		}
	}
	return out
}
func percent(v volumeMap) int {
	for _, c := range v {
		var p int
		fmt.Sscanf(c.ValuePercent, "%d%%", &p)
		return p
	}
	return 0
}
func bar(p, w int) string {
	fill := min(w, max(0, p*w/100))
	return "[" + strings.Repeat("█", fill) + strings.Repeat("░", w-fill) + "]"
}
func appName(s stream) string {
	for _, k := range []string{"application.name", "media.name", "node.name"} {
		if s.Properties[k] != "" {
			return s.Properties[k]
		}
	}
	return fmt.Sprintf("stream %d", s.Index)
}

/* Shared helpers */
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
func short(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
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

var _ = context.Background

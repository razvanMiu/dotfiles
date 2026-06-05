package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
type state struct {
	Info                  info
	Sinks, Sources        []device
	Inputs, SourceOutputs []stream
}
type tickMsg struct{}

type model struct {
	st              state
	section, cursor int
	width, height   int
	err             string
}

const (
	base     = "#24273a"
	surface0 = "#363a4f"
	text     = "#cad3f5"
	subtext  = "#a5adcb"
	mauve    = "#c6a0f6"
	green    = "#a6da95"
	red      = "#ed8796"
	yellow   = "#eed49f"
)

var (
	page  = lipgloss.NewStyle().Background(lipgloss.Color(base)).Foreground(lipgloss.Color(text)).Padding(1, 2)
	title = lipgloss.NewStyle().Foreground(lipgloss.Color(mauve)).Bold(true)
	muted = lipgloss.NewStyle().Foreground(lipgloss.Color(subtext))
	sel   = lipgloss.NewStyle().Foreground(lipgloss.Color(base)).Background(lipgloss.Color(mauve))
	ok    = lipgloss.NewStyle().Foreground(lipgloss.Color(green))
	warn  = lipgloss.NewStyle().Foreground(lipgloss.Color(yellow))
	bad   = lipgloss.NewStyle().Foreground(lipgloss.Color(red))
)

func main() {
	m := model{}
	m.refresh()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func (m model) Init() tea.Cmd { return tick() }
func tick() tea.Cmd           { return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return tickMsg{} }) }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.refresh()
		return m, tick()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
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
	return m, nil
}
func (m model) View() string {
	w := max(58, m.width-4)
	var b strings.Builder
	writeLine(&b, center(title.Render("Audio"), w))
	writeLine(&b, center(muted.Render("Tab section · j/k row · h/l vol · m mute · Enter default · r/q"), w))
	writeBlankLine(&b)
	b.WriteString(m.renderSection(0, "Outputs", m.renderDevices(m.st.Sinks, m.st.Info.DefaultSinkName, true), w))
	b.WriteString(m.renderSection(1, "Inputs", m.renderDevices(m.st.Sources, m.st.Info.DefaultSourceName, false), w))
	b.WriteString(m.renderSection(2, "Apps", m.renderStreams(m.st.Inputs), w))
	b.WriteString(m.renderSection(3, "Recording", m.renderStreams(m.st.SourceOutputs), w))
	if m.err != "" {
		writeBlankLine(&b)
		b.WriteString(bad.Render(m.err))
	}
	return page.Width(w).Render(b.String())
}
func (m model) renderSection(i int, name string, lines []string, w int) string {
	head := name
	if m.section == i {
		head = "▶ " + head
	} else {
		head = "  " + head
	}
	var b strings.Builder
	b.WriteString(title.Render(head))
	b.WriteString(muted.Render(fmt.Sprintf("  %d", len(lines))))
	b.WriteByte('\n')
	if len(lines) == 0 {
		b.WriteString("  ")
		writeLine(&b, muted.Render("No active streams"))
		writeBlankLine(&b)
		return b.String()
	}
	for n, l := range lines {
		if m.section == i && m.cursor == n {
			l = sel.Render(l)
		}
		b.WriteString("  ")
		writeLine(&b, l)
	}
	writeBlankLine(&b)
	return b.String()
}
func (m model) renderDevices(ds []device, def string, output bool) []string {
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
		out = append(out, fmt.Sprintf("%s %-30s %3d%% %s%s", mark, short(d.Description, 30), percent(d.Volume), bar(percent(d.Volume), 16), mute))
	}
	return out
}
func (m model) renderStreams(ss []stream) []string {
	var out []string
	sort.Slice(ss, func(i, j int) bool { return appName(ss[i]) < appName(ss[j]) })
	for _, s := range ss {
		mute := ""
		if s.Mute {
			mute = " " + bad.Render("muted")
		}
		out = append(out, fmt.Sprintf("%-32s %3d%% %s%s", short(appName(s), 32), percent(s.Volume), bar(percent(s.Volume), 16), mute))
	}
	return out
}
func (m *model) refresh() {
	st, err := load()
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
func load() (state, error) {
	var st state
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
func (m model) sectionLen() int {
	switch m.section {
	case 0:
		return len(m.renderDevices(m.st.Sinks, m.st.Info.DefaultSinkName, true))
	case 1:
		return len(m.renderDevices(m.st.Sources, m.st.Info.DefaultSourceName, false))
	case 2:
		return len(m.st.Inputs)
	case 3:
		return len(m.st.SourceOutputs)
	}
	return 0
}
func (m model) selectedID() string {
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
func (m *model) adjust(delta int) {
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
func (m *model) toggleMute() {
	id := m.selectedID()
	if id == "" {
		return
	}
	cmds := [][]string{{"set-sink-mute"}, {"set-source-mute"}, {"set-sink-input-mute"}, {"set-source-output-mute"}}
	exec.Command("pactl", cmds[m.section][0], id, "toggle").Run()
	m.refresh()
	refreshStatusBar()
}
func (m *model) activate() {
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
	cmd := dwmStatusPath()
	if cmd == "" {
		return
	}
	exec.Command(cmd, "--refresh", "volume").Run()
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
func short(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}
func center(s string, w int) string {
	return lipgloss.NewStyle().Width(w).Align(lipgloss.Center).Render(s)
}
func writeLine(b *strings.Builder, s string) {
	b.WriteString(s)
	b.WriteByte('\n')
}
func writeBlankLine(b *strings.Builder) {
	b.WriteByte('\n')
}

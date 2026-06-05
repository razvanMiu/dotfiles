package panels

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
type audioState struct {
	Info                  info
	Sinks, Sources        []device
	Inputs, SourceOutputs []stream
}
type Audio struct {
	st              audioState
	section, cursor int
	width, height   int
	err             string
}

func NewAudio() Audio                      { m := Audio{}; m.Refresh(); return m }
func (m *Audio) SetSize(width, height int) { m.width, m.height = width, height }
func (m Audio) WidthHint() int             { return 76 }
func (m Audio) HeightHint() int            { return 32 }
func (m *Audio) Update(msg tea.Msg) {
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
			m.Refresh()
		}
	}
}
func (m Audio) View(w, h int) string {
	lines := []string{
		center(title.Render("Audio"), w),
		center(muted.Render("Tab section · j/k row · h/l vol · m mute · Enter default · r/q"), w),
		"",
	}
	sections := []struct {
		index int
		name  string
		lines []string
	}{
		{0, "Outputs", m.renderDevices(m.st.Sinks, m.st.Info.DefaultSinkName, true)},
		{1, "Inputs", m.renderDevices(m.st.Sources, m.st.Info.DefaultSourceName, false)},
		{2, "Apps", m.renderStreams(m.st.Inputs)},
		{3, "Recording", m.renderStreams(m.st.SourceOutputs)},
	}
	remaining := max(0, h-len(lines))
	for i, section := range sections {
		sectionsLeft := len(sections) - i
		budget := max(1, remaining-sectionsLeft+1)
		rendered := m.renderSection(section.index, section.name, section.lines, budget)
		lines = append(lines, rendered...)
		remaining = max(0, h-len(lines))
	}
	if m.err != "" && len(lines) < h {
		lines = append(lines, bad.Render(m.err))
	}
	return fitLines(lines, w, h)
}
func (m Audio) renderSection(i int, name string, rows []string, budget int) []string {
	head := "  " + name
	if m.section == i {
		head = "▶ " + name
	}
	lines := []string{title.Render(head) + muted.Render(fmt.Sprintf("  %d", len(rows)))}
	if budget <= 1 {
		return lines
	}
	if len(rows) == 0 {
		return append(lines, "  "+muted.Render("No active streams"))
	}
	rowBudget := max(0, budget-1)
	start := 0
	if m.section == i && m.cursor >= rowBudget && rowBudget > 0 {
		start = m.cursor - rowBudget + 1
	}
	end := min(len(rows), start+rowBudget)
	for n, l := range rows[start:end] {
		idx := start + n
		if m.section == i && m.cursor == idx {
			l = selected.Render(l)
		}
		lines = append(lines, "  "+l)
	}
	if end < len(rows) && len(lines) > 1 {
		lines[len(lines)-1] = "  " + muted.Render("↓ more")
	}
	return lines
}
func (m Audio) renderDevices(ds []device, def string, output bool) []string {
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
func (m Audio) renderStreams(ss []stream) []string {
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
func (m *Audio) Refresh() {
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
func (m Audio) sectionLen() int {
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
func (m Audio) selectedID() string {
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
func (m *Audio) adjust(delta int) {
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
	m.Refresh()
	refreshStatusBar()
}
func (m *Audio) toggleMute() {
	id := m.selectedID()
	if id == "" {
		return
	}
	cmds := []string{"set-sink-mute", "set-source-mute", "set-sink-input-mute", "set-source-output-mute"}
	exec.Command("pactl", cmds[m.section], id, "toggle").Run()
	m.Refresh()
	refreshStatusBar()
}
func (m *Audio) activate() {
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
	m.Refresh()
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

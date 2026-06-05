package main

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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

func readFile(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}

func stateDir() string {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return filepath.Join(v, "dwm")
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "state", "dwm")
}

func socketPath() string  { return filepath.Join(stateDir(), "widgets.sock") }
func lastTabPath() string { return filepath.Join(stateDir(), "last-widget-tab") }

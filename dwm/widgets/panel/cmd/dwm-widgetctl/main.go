package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "calendar" && os.Args[1] != "audio") {
		fmt.Fprintln(os.Stderr, "usage: dwm-widgetctl <calendar|audio>")
		os.Exit(2)
	}
	tab := os.Args[1]
	_ = os.MkdirAll(stateDir(), 0755)
	_ = os.WriteFile(lastTabPath(), []byte(tab), 0644)
	addr, err := net.ResolveUnixAddr("unixgram", socketPath())
	if err != nil {
		return
	}
	conn, err := net.DialUnix("unixgram", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()
	_ = conn.SetWriteDeadline(time.Now().Add(200 * time.Millisecond))
	_, _ = conn.Write([]byte(tab))
}

func stateDir() string {
	if v := os.Getenv("XDG_STATE_HOME"); strings.TrimSpace(v) != "" {
		return filepath.Join(v, "dwm")
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "state", "dwm")
}
func socketPath() string  { return filepath.Join(stateDir(), "widgets.sock") }
func lastTabPath() string { return filepath.Join(stateDir(), "last-widget-tab") }

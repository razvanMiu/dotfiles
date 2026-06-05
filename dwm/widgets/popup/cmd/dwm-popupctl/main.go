package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dwm-popupctl show|toggle [calendar|audio] | hide | quit")
		os.Exit(2)
	}
	conn, err := net.Dial("unix", socketPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "dwm-popupctl: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	_, _ = fmt.Fprintln(conn, strings.Join(os.Args[1:], " "))
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil && n == 0 {
		fmt.Fprintf(os.Stderr, "dwm-popupctl: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(buf[:n]))
}

func socketPath() string {
	if state := os.Getenv("XDG_STATE_HOME"); state != "" {
		return filepath.Join(state, "dwm", "popup.sock")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(os.TempDir(), "dwm-popup.sock")
	}
	return filepath.Join(home, ".local", "state", "dwm", "popup.sock")
}

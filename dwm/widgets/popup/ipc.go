package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type command struct {
	action string
	panel  string
	resp   chan string
}

func parseCommand(line string) (command, error) {
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(line)))
	if len(fields) == 0 {
		return command{}, errors.New("empty command")
	}
	cmd := command{action: fields[0]}
	switch cmd.action {
	case "hide", "quit":
		if len(fields) != 1 {
			return command{}, fmt.Errorf("%s takes no arguments", cmd.action)
		}
	case "show", "toggle":
		if len(fields) > 2 {
			return command{}, fmt.Errorf("%s takes at most one panel", cmd.action)
		}
		cmd.panel = "calendar"
		if len(fields) == 2 {
			cmd.panel = normalizePanel(fields[1])
		}
	default:
		return command{}, fmt.Errorf("unknown command %q", cmd.action)
	}
	return cmd, nil
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

func startIPC(commands chan<- command) (func(), error) {
	path := socketPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if conn, err := net.Dial("unix", path); err == nil {
		_ = conn.Close()
		return nil, fmt.Errorf("popup IPC already active at %s", path)
	}
	_ = os.Remove(path)
	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handleConn(conn, commands)
		}
	}()
	return func() {
		_ = ln.Close()
		_ = os.Remove(path)
	}, nil
}

func handleConn(conn net.Conn, commands chan<- command) {
	defer conn.Close()
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil && strings.TrimSpace(line) == "" {
		_, _ = fmt.Fprintf(conn, "error: %v\n", err)
		return
	}
	cmd, err := parseCommand(line)
	if err != nil {
		_, _ = fmt.Fprintf(conn, "error: %v\n", err)
		return
	}
	cmd.resp = make(chan string, 1)
	commands <- cmd
	_, _ = fmt.Fprintln(conn, <-cmd.resp)
}

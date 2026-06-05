package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"dwm-status/internal/segments"
	"dwm-status/internal/status"
)

func main() {
	tuneRuntime()
	stateDir := stateDir()
	socket := filepath.Join(stateDir, "status.sock")

	if len(os.Args) > 1 {
		runClient(socket, stateDir, os.Args[1:])
		return
	}
	runDaemon(stateDir, socket)
}

func runDaemon(stateDir, socket string) {
	logger := log.New(os.Stderr, "dwm-status: ", log.LstdFlags)
	d := status.NewDaemon(segments.New(stateDir), status.XRootPublisher{StateDir: stateDir})
	d.Log = logger
	ctx := context.Background()
	if err := d.Refresh(ctx, "all"); err != nil {
		logger.Printf("initial render: %v", err)
	}
	go func() {
		if err := status.Serve(ctx, socket, d); err != nil {
			logger.Printf("ipc: %v", err)
			os.Exit(1)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for now := range ticker.C {
		if err := d.RefreshDue(ctx, now); err != nil {
			logger.Printf("render: %v", err)
		}
	}
}

func runClient(socket, stateDir string, args []string) {
	switch {
	case len(args) == 1 && args[0] == "--health":
		if err := status.Health(socket, time.Second); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case len(args) == 2 && args[0] == "--refresh":
		if err := status.Send(socket, "refresh "+args[1], time.Second); err != nil {
			if args[1] == "all" {
				oneShot(stateDir)
				return
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "usage: dwm-status [--health | --refresh <segment|all>]")
		os.Exit(2)
	}
}

func oneShot(stateDir string) {
	d := status.NewDaemon(segments.New(stateDir), status.XRootPublisher{StateDir: stateDir})
	if err := d.Refresh(context.Background(), "all"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func tuneRuntime() {
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(25)
	debug.SetMemoryLimit(16 * 1024 * 1024)
}

func stateDir() string {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return filepath.Join(v, "dwm")
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "state", "dwm")
}

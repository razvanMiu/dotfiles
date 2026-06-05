package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDwmStatusPathUsesLocalBinWhenPathDoesNotContainIt(t *testing.T) {
	home := t.TempDir()
	bin := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(bin, 0755); err != nil {
		t.Fatal(err)
	}
	status := filepath.Join(bin, "dwm-status")
	if err := os.WriteFile(status, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("PATH", "/usr/bin:/bin")
	if got := dwmStatusPath(); got != status {
		t.Fatalf("dwmStatusPath() = %q, want %q", got, status)
	}
}

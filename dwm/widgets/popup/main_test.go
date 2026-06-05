package main

import (
	"testing"
	"time"
)

func TestCalendarGridStartsOnMonday(t *testing.T) {
	grid := calendarGrid(time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC))
	if got := grid[0][0]; got.Day() != 1 || got.Month() != time.June {
		t.Fatalf("grid[0][0] = %v, want 2026-06-01", got)
	}
	if got := grid[0][6]; got.Day() != 7 || got.Weekday() != time.Sunday {
		t.Fatalf("grid[0][6] = %v, want Sunday 2026-06-07", got)
	}
}

func TestCalendarGridIncludesPreviousMonthPadding(t *testing.T) {
	grid := calendarGrid(time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC))
	if got := grid[0][0]; got.Day() != 27 || got.Month() != time.July {
		t.Fatalf("grid[0][0] = %v, want 2026-07-27", got)
	}
	if got := grid[0][5]; got.Day() != 1 || got.Month() != time.August {
		t.Fatalf("grid[0][5] = %v, want 2026-08-01", got)
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input  string
		action string
		panel  string
	}{
		{"show", "show", "calendar"},
		{"show audio", "show", "audio"},
		{"toggle calendar", "toggle", "calendar"},
		{"hide", "hide", ""},
		{"quit", "quit", ""},
	}
	for _, tt := range tests {
		cmd, err := parseCommand(tt.input)
		if err != nil {
			t.Fatalf("parseCommand(%q): %v", tt.input, err)
		}
		if cmd.action != tt.action || cmd.panel != tt.panel {
			t.Fatalf("parseCommand(%q) = (%q, %q), want (%q, %q)", tt.input, cmd.action, cmd.panel, tt.action, tt.panel)
		}
	}
}

func TestParseCommandRejectsBadInput(t *testing.T) {
	for _, input := range []string{"", "wat", "hide audio", "show audio extra"} {
		if _, err := parseCommand(input); err == nil {
			t.Fatalf("parseCommand(%q) succeeded, want error", input)
		}
	}
}

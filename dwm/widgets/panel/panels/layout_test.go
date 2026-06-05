package panels

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func assertFits(t *testing.T, name, view string, width, height int) {
	t.Helper()
	lines := strings.Split(view, "\n")
	if len(lines) != height {
		t.Fatalf("%s height = %d, want %d", name, len(lines), height)
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got > width {
			t.Fatalf("%s line %d width = %d, want <= %d: %q", name, i, got, width, line)
		}
	}
}

func TestCalendarViewFitsRequestedCells(t *testing.T) {
	cal := NewCalendar()
	assertFits(t, "calendar", cal.View(44, 28), 44, 28)
}

func TestAudioViewFitsRequestedCells(t *testing.T) {
	audio := Audio{}
	assertFits(t, "audio", audio.View(76, 30), 76, 30)
}

package segments

import (
	"context"
	"os/exec"
	"time"

	"dwm-status/internal/status"
)

type commandRunner func(context.Context, string, ...string) ([]byte, error)

func command(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).Output()
}

func New(stateDir string) []status.Segment {
	return []status.Segment{
		NewCPU(command, 2*time.Second),
		NewGPU(command, 2*time.Second),
		NewRAM(2 * time.Second),
		NewVolume(command, 2*time.Second),
		NewTime(30 * time.Second),
	}
}

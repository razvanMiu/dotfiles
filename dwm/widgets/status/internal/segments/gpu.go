package segments

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dwm-status/internal/status"
)

type GPU struct {
	run      commandRunner
	interval time.Duration
}

func NewGPU(run commandRunner, interval time.Duration) *GPU {
	return &GPU{run: run, interval: interval}
}
func (g *GPU) Name() string            { return "gpu" }
func (g *GPU) Interval() time.Duration { return g.interval }
func (g *GPU) Update(ctx context.Context) (status.Value, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	out, err := g.run(ctx, "nvidia-smi",
		"--query-gpu=utilization.gpu,temperature.gpu,memory.used,memory.total",
		"--format=csv,noheader,nounits")
	if err != nil {
		return status.Value{Text: "󰢮 --% --°C --/--G"}, nil
	}
	return status.Value{Text: renderGPU(string(out))}, nil
}

func renderGPU(out string) string {
	line := strings.TrimSpace(strings.Split(strings.TrimSpace(out), "\n")[0])
	parts := strings.Split(line, ",")
	if len(parts) != 4 {
		return "󰢮 --% --°C --/--G"
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	used, usedErr := strconv.ParseFloat(parts[2], 64)
	total, totalErr := strconv.ParseFloat(parts[3], 64)
	if parts[0] == "" || parts[1] == "" || usedErr != nil || totalErr != nil || total <= 0 {
		return "󰢮 --% --°C --/--G"
	}
	return fmt.Sprintf("󰢮 %s%% %s°C %.1f/%.0fG", parts[0], parts[1], used/1024, total/1024)
}

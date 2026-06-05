package segments

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dwm-status/internal/status"
)

type Volume struct {
	run      commandRunner
	interval time.Duration
}

func NewVolume(run commandRunner, interval time.Duration) *Volume {
	return &Volume{run: run, interval: interval}
}
func (v *Volume) Name() string            { return "volume" }
func (v *Volume) Interval() time.Duration { return v.interval }
func (v *Volume) Update(ctx context.Context) (status.Value, error) {
	out, err := v.run(ctx, "wpctl", "get-volume", "@DEFAULT_AUDIO_SINK@")
	if err != nil {
		return status.Value{Text: "󰕾 --%"}, nil
	}
	fields := strings.Fields(string(out))
	if len(fields) < 2 {
		return status.Value{Text: "󰕾 --%"}, nil
	}
	level, _ := strconv.ParseFloat(fields[1], 64)
	icon := "󰕾"
	if strings.Contains(string(out), "[MUTED]") {
		icon = "󰖁"
	}
	return status.Value{Text: fmt.Sprintf("%s %.0f%%", icon, level*100)}, nil
}

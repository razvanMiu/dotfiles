package segments

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"dwm-status/internal/status"
)

type CPU struct {
	run      commandRunner
	interval time.Duration
	prev     cpuSample
}

type cpuSample struct{ total, idle uint64 }

func NewCPU(run commandRunner, interval time.Duration) *CPU {
	return &CPU{run: run, interval: interval}
}
func (c *CPU) Name() string            { return "cpu" }
func (c *CPU) Interval() time.Duration { return c.interval }
func (c *CPU) Update(ctx context.Context) (status.Value, error) {
	sample := readCPU()
	usage := 0
	if c.prev.total != 0 && sample.total > c.prev.total {
		dt := sample.total - c.prev.total
		di := sample.idle - c.prev.idle
		usage = int((100 * (dt - di)) / dt)
	}
	c.prev = sample
	return status.Value{Text: fmt.Sprintf("󰘚 %d%% %s°C", usage, c.temp(ctx))}, nil
}

func readCPU() cpuSample {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return cpuSample{}
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	if !s.Scan() {
		return cpuSample{}
	}
	fields := strings.Fields(s.Text())
	var vals []uint64
	for _, field := range fields[1:] {
		v, _ := strconv.ParseUint(field, 10, 64)
		vals = append(vals, v)
	}
	if len(vals) < 7 {
		return cpuSample{}
	}
	idle := vals[3] + vals[4]
	total := uint64(0)
	for _, v := range vals[:7] {
		total += v
	}
	return cpuSample{total: total, idle: idle}
}

func (c *CPU) temp(ctx context.Context) string {
	out, err := c.run(ctx, "sensors")
	if err != nil {
		return "--"
	}
	maxTemp := math.NaN()
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if !(strings.HasPrefix(line, "Package id 0:") || strings.HasPrefix(line, "Core ")) {
			continue
		}
		for _, field := range strings.Fields(line) {
			if strings.HasPrefix(field, "+") && strings.HasSuffix(field, "°C") {
				v, err := strconv.ParseFloat(strings.Trim(field, "+°C"), 64)
				if err == nil && (math.IsNaN(maxTemp) || v > maxTemp) {
					maxTemp = v
				}
			}
		}
	}
	if math.IsNaN(maxTemp) {
		return "--"
	}
	return fmt.Sprintf("%.0f", maxTemp)
}

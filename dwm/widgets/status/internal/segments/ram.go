package segments

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"dwm-status/internal/status"
)

type RAM struct{ interval time.Duration }

func NewRAM(interval time.Duration) *RAM { return &RAM{interval: interval} }
func (r *RAM) Name() string              { return "ram" }
func (r *RAM) Interval() time.Duration   { return r.interval }
func (r *RAM) Update(context.Context) (status.Value, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return status.Value{Text: "󰍛 --/--G"}, nil
	}
	defer f.Close()
	return status.Value{Text: renderRAM(f)}, nil
}

func renderRAM(r io.Reader) string {
	/*
		Use MemTotal - MemAvailable for real system RAM pressure.
		Reserved but untouched virtual address space, like Go VSZ/VIRT, is not
		resident RAM and is intentionally excluded by this /proc/meminfo path.
	*/
	total, available := readMeminfo(r)
	if total <= 0 || available <= 0 {
		return "󰍛 --/--G"
	}
	used := (total - available) / 1024 / 1024
	return fmt.Sprintf("󰍛 %.1f/%.0fG", used, total/1024/1024)
}

func readMeminfo(r io.Reader) (total, available float64) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		p := strings.Fields(sc.Text())
		if len(p) < 2 {
			continue
		}
		v, _ := strconv.ParseFloat(p[1], 64)
		switch strings.TrimSuffix(p[0], ":") {
		case "MemTotal":
			total = v
		case "MemAvailable":
			available = v
		}
		if total > 0 && available > 0 {
			break
		}
	}
	return total, available
}

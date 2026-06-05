package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const DefaultSeparator = "\x1f"

type Publisher interface {
	Publish(text string) error
}

type XRootPublisher struct {
	StateDir string
}

func (p XRootPublisher) Publish(text string) error {
	if err := os.MkdirAll(p.StateDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(p.StateDir, "status-current"), []byte(text), 0644); err != nil {
		return err
	}
	return exec.Command("xsetroot", "-name", text).Run()
}

type Formatter struct {
	Separator string
}

func (f Formatter) Format(values map[string]Value, segments []Segment) string {
	sep := f.Separator
	if sep == "" {
		sep = DefaultSeparator
	}
	var b strings.Builder
	for i, seg := range segments {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(values[seg.Name()].Text)
	}
	return b.String()
}

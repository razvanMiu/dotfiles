package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"dwm-status/internal/segments"
	"dwm-status/internal/status"
)

type fakeSegment struct {
	name     string
	interval time.Duration
	value    string
	updates  int
}

func (f *fakeSegment) Name() string            { return f.name }
func (f *fakeSegment) Interval() time.Duration { return f.interval }
func (f *fakeSegment) Update(context.Context) (status.Value, error) {
	f.updates++
	return status.Value{Text: f.value}, nil
}

type capturePublisher struct{ text string }

func (p *capturePublisher) Publish(text string) error { p.text = text; return nil }

func TestTargetedRefreshOnlyUpdatesRequestedSegment(t *testing.T) {
	cpu := &fakeSegment{name: "cpu", interval: time.Minute, value: "cpu-old"}
	vol := &fakeSegment{name: "volume", interval: time.Minute, value: "vol-old"}
	publisher := &capturePublisher{}
	d := status.NewDaemon([]status.Segment{cpu, vol}, publisher)
	if err := d.Refresh(context.Background(), "all"); err != nil {
		t.Fatal(err)
	}
	if publisher.text != "cpu-old\x1fvol-old" {
		t.Fatalf("initial publish = %q", publisher.text)
	}

	cpu.value = "cpu-new"
	vol.value = "vol-new"
	if err := d.Refresh(context.Background(), "volume"); err != nil {
		t.Fatal(err)
	}
	if publisher.text != "cpu-old\x1fvol-new" {
		t.Fatalf("targeted publish = %q", publisher.text)
	}
	if cpu.updates != 1 {
		t.Fatalf("cpu updates = %d, want 1", cpu.updates)
	}
	if vol.updates != 2 {
		t.Fatalf("volume updates = %d, want 2", vol.updates)
	}
}

func TestRefreshDueHonorsSegmentIntervals(t *testing.T) {
	fast := &fakeSegment{name: "fast", interval: time.Nanosecond, value: "f1"}
	slow := &fakeSegment{name: "slow", interval: time.Hour, value: "s1"}
	publisher := &capturePublisher{}
	d := status.NewDaemon([]status.Segment{fast, slow}, publisher)
	if err := d.Refresh(context.Background(), "all"); err != nil {
		t.Fatal(err)
	}
	fast.value = "f2"
	slow.value = "s2"
	time.Sleep(time.Millisecond)
	if err := d.RefreshDue(context.Background(), time.Now()); err != nil {
		t.Fatal(err)
	}
	if publisher.text != "f2\x1fs1" {
		t.Fatalf("due publish = %q", publisher.text)
	}
}

func TestGPUFormatsNvidiaSMIOutput(t *testing.T) {
	gpu := segments.NewGPU(func(context.Context, string, ...string) ([]byte, error) {
		return []byte("28, 36, 614, 12282\n"), nil
	}, time.Second)
	value, err := gpu.Update(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "󰢮 28% 36°C 0.6/12G"
	if value.Text != want {
		t.Fatalf("gpu text = %q, want %q", value.Text, want)
	}
}

func TestGPUFallsBackWhenNvidiaSMIFails(t *testing.T) {
	gpu := segments.NewGPU(func(context.Context, string, ...string) ([]byte, error) {
		return nil, errors.New("nvidia-smi unavailable")
	}, time.Second)
	value, err := gpu.Update(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "󰢮 --% --°C --/--G"
	if value.Text != want {
		t.Fatalf("gpu text = %q, want %q", value.Text, want)
	}
}

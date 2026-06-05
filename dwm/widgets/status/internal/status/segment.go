package status

import (
	"context"
	"time"
)

type Value struct {
	Text string
}

type Segment interface {
	Name() string
	Interval() time.Duration
	Update(context.Context) (Value, error)
}

type SegmentFunc struct {
	SegmentName     string
	RefreshInterval time.Duration
	UpdateFunc      func(context.Context) (Value, error)
}

func (s SegmentFunc) Name() string                              { return s.SegmentName }
func (s SegmentFunc) Interval() time.Duration                   { return s.RefreshInterval }
func (s SegmentFunc) Update(ctx context.Context) (Value, error) { return s.UpdateFunc(ctx) }

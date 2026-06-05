package status

import (
	"context"
	"log"
	"time"
)

type Daemon struct {
	Segments  []Segment
	Publisher Publisher
	Formatter Formatter
	Values    map[string]Value
	next      map[string]time.Time
	Log       *log.Logger
}

func NewDaemon(segments []Segment, publisher Publisher) *Daemon {
	return &Daemon{
		Segments:  segments,
		Publisher: publisher,
		Formatter: Formatter{Separator: DefaultSeparator},
		Values:    map[string]Value{},
		next:      map[string]time.Time{},
	}
}

func (d *Daemon) Refresh(ctx context.Context, name string) error {
	if name == "" || name == "all" {
		for _, seg := range d.Segments {
			d.update(ctx, seg)
		}
	} else {
		for _, seg := range d.Segments {
			if seg.Name() == name {
				d.update(ctx, seg)
				break
			}
		}
	}
	return d.publish()
}

func (d *Daemon) RefreshDue(ctx context.Context, now time.Time) error {
	changed := false
	for _, seg := range d.Segments {
		if due, ok := d.next[seg.Name()]; !ok || !now.Before(due) {
			d.update(ctx, seg)
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return d.publish()
}

func (d *Daemon) publish() error {
	return d.Publisher.Publish(d.Formatter.Format(d.Values, d.Segments))
}

func (d *Daemon) update(ctx context.Context, seg Segment) {
	value, err := seg.Update(ctx)
	if err != nil {
		if d.Log != nil {
			d.Log.Printf("segment %s: %v", seg.Name(), err)
		}
		if _, ok := d.Values[seg.Name()]; ok {
			return
		}
		value = Value{Text: seg.Name() + " --"}
	}
	d.Values[seg.Name()] = value
	d.next[seg.Name()] = time.Now().Add(seg.Interval())
}

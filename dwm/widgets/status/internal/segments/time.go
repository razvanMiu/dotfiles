package segments

import (
	"context"
	"time"

	"dwm-status/internal/status"
)

type Time struct{ interval time.Duration }

func NewTime(interval time.Duration) *Time { return &Time{interval: interval} }
func (t *Time) Name() string               { return "time" }
func (t *Time) Interval() time.Duration    { return t.interval }
func (t *Time) Update(context.Context) (status.Value, error) {
	return status.Value{Text: "󰥔 " + time.Now().Format("Mon 02 Jan 15:04")}, nil
}

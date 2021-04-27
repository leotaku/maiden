package diary

import (
	"strings"
	"time"

	"github.com/leotaku/maiden/timeutil"
)

type Entry struct {
	Datetime    timeutil.Datetime
	Duration    time.Duration
	Description string
}

func (e Entry) Start(loc *time.Location) time.Time {
	return e.Datetime.First(loc)
}

func (e Entry) End(loc *time.Location) time.Time {
	return e.Datetime.First(loc).Add(e.Duration)
}

func (e Entry) RRule(loc *time.Location) string {
	return e.Datetime.Date.RRule(loc)
}

func (e Entry) Format(o timeutil.DateOrder) string {
	parts := make([]string, 0)
	parts = append(parts, e.Datetime.Format(o))
	if e.Duration != time.Hour*24 {
		parts = append(parts, timeutil.FormatDuration(e.Duration))
	}
	parts = append(parts, e.Description)

	return strings.Join(parts, " ")
}

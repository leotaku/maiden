package diary

import (
	"time"

	"github.com/leotaku/maiden/timeutil"
)

var (
	ISO      = timeutil.ISO
	American = timeutil.American
	European = timeutil.European
)

type Diary []Entry

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

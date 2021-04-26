package timeutil

import (
	"time"
)

type Datetime struct {
	Date Date
	Time Time
}

type Date struct {
	Year     *int
	Month    *time.Month
	Monthday *int
	Weekday  *time.Weekday
}

type Time struct {
	Hour   *int
	Minute *int
}

func (d Datetime) First(loc *time.Location) time.Time {
	year := time.Now().Year()
	month := time.January
	monthday := 1
	hour := 0
	minute := 0

	if d.Date.Year != nil {
		year = *d.Date.Year
	}
	if d.Date.Month != nil {
		month = *d.Date.Month
	}
	if d.Date.Monthday != nil {
		monthday = *d.Date.Monthday
	}
	if d.Time.Hour != nil {
		hour = *d.Time.Hour
	}
	if d.Time.Minute != nil {
		minute = *d.Time.Minute
	}

	date := time.Date(year, month, monthday, hour, minute, 0, 0, loc)
	if d.Date.Weekday != nil {
		for date.Weekday() != *d.Date.Weekday {
			date = date.Add(time.Hour * 24)
		}
	}

	return date
}

package timeutil

import (
	"fmt"
	"strings"
	"time"
)

const (
	DateTimeFormat = "20060102T150405Z"
	DateFormat     = "20060102"
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

func (d Date) RRule(loc *time.Location) string {
	if d.Year != nil && d.Month != nil && d.Monthday != nil {
		return ""
	}

	rules := make([]string, 0)
	rules = append(rules, "FREQ=DAILY")

	if d.Year != nil {
		date := time.Date(*d.Year+1, 0, 0, 0, 0, 0, 0, loc)
		s := fmt.Sprintf("UNTIL=%v", date.Format(DateTimeFormat))
		rules = append(rules, s)
	}
	if d.Month != nil {
		s := fmt.Sprintf("BYMONTH=%v", int(*d.Month))
		rules = append(rules, s)
	}
	if d.Monthday != nil {
		s := fmt.Sprintf("BYMONTHDAY=%v", *d.Monthday)
		rules = append(rules, s)
		rules[0] = "FREQ=MONTHLY"
	}
	if d.Weekday != nil {
		wday := strings.ToUpper(d.Weekday.String()[:2])
		s := fmt.Sprintf("BYWEEKDAY=%v", wday)
		rules = append(rules, s)
		rules[0] = "FREQ=WEEKLY"
	}

	return strings.Join(rules, ";")
}

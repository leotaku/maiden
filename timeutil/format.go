package timeutil

import (
	"fmt"
	"strings"
	"time"
)

const (
	DateTimeFormat = "20060102T150405"
	DateFormat     = "20060102"
)

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

func (d Datetime) Format(o DateOrder) string {
	date := d.Date.Format(o)
	time := d.Time.String()

	if date == "" && time == "" {
		return ""
	} else if date == "" {
		return fmt.Sprintf(" %v", time)
	} else {
		return fmt.Sprintf("%v %v", date, time)
	}
}

func (d Date) Format(o DateOrder) string {
	year := ""
	if d.Year != nil {
		year = fmt.Sprintf("%04d", *d.Year)
	}
	month := ""
	if d.Month != nil {
		month = fmt.Sprintf("%02d", *d.Month)
	}
	monthday := ""
	if d.Monthday != nil {
		monthday = fmt.Sprintf("%02d", *d.Monthday)
	}

	ordered := make([]string, 0)
	left, middle, right := o(year, month, monthday)
	if left != "" {
		ordered = append(ordered, left)
	}
	if middle != "" {
		ordered = append(ordered, middle)
	}
	if right != "" {
		ordered = append(ordered, right)
	}

	parts := make([]string, 0)
	if ymd := strings.Join(ordered, "-"); ymd != "" {
		parts = append(parts, ymd)
	}
	if d.Weekday != nil {
		parts = append(parts, d.Weekday.String())
	}

	return strings.Join(parts, " ")
}

func (t Time) String() string {
	hour := 0
	minute := 0
	if t.Hour != nil {
		hour = *t.Hour
	}
	if t.Minute != nil {
		minute = *t.Minute
	}

	return fmt.Sprintf("%02d:%02d", hour, minute)
}

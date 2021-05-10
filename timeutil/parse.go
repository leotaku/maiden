package timeutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// If the smaller (year>month>day) provided value is the empty string,
// the other values are shifted to replace it.
func ParseDate(left, middle, right string, o DateOrder) (*Date, error) {
	year, month, day := o(left, middle, right)

	if day == "" {
		day = month
		month = year
		year = ""
	}
	if day == "" {
		day = month
		month = ""
	}

	return parseFullISODate(year, month, day)
}

func parseFullISODate(year, month, day string) (*Date, error) {
	y, err := parseRange(year, 0, math.MaxInt64)
	if year != "" && err != nil {
		return nil, fmt.Errorf("year: %w", err)
	}
	m, err := ParseMonth(month)
	if month != "" && err != nil {
		return nil, fmt.Errorf("month: %w", err)
	}
	d, err := parseRange(day, 1, 31)
	if day != "" && err != nil {
		return nil, fmt.Errorf("day: %w", err)
	}

	return &Date{
		Year:     y,
		Month:    m,
		Monthday: d,
	}, nil
}

func ParseTime(hour, minute string) (*Time, error) {
	h, err := parseRange(hour, 0, 24)
	if hour != "" && err != nil {
		return nil, fmt.Errorf("hour: %w", err)
	}
	m, err := parseRange(minute, 0, 60)
	if hour != "" && err != nil {
		return nil, fmt.Errorf("minute: %w", err)
	}

	return &Time{
		Hour:   h,
		Minute: m,
	}, nil
}

func ParseMonth(month string) (*time.Month, error) {
	if m, err := parseRange(month, 1, 12); err == nil {
		return (*time.Month)(m), nil
	} else {
		for i := time.January; i <= time.December; i++ {
			if strings.HasPrefix(month, i.String()[0:3]) {
				return &i, nil
			}
		}
	}

	return nil, fmt.Errorf("not a month: %v", month)
}

func ParseWeekday(wday string) (*time.Weekday, error) {
	for i := time.Sunday; i <= time.Saturday; i++ {
		if strings.HasPrefix(wday, i.String()[0:2]) {
			return &i, nil
		}
	}

	return nil, fmt.Errorf("not a weekday: %v", wday)
}

func parseRange(s string, min, max int) (*int, error) {
	it, err := strconv.Atoi(s)
	if s == "*" {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else if it < min || it > max {
		return nil, fmt.Errorf("out of range: %v", it)
	} else {
		return &it, nil
	}
}

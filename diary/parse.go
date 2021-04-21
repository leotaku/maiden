package diary

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/leotaku/maiden/timeutil"
)

type Parser struct {
	reader   *bufio.Reader
	lastDate *timeutil.Date
	order    timeutil.DateOrder
}

func NewParser(r io.Reader, order timeutil.DateOrder) *Parser {
	return &Parser{
		reader:   bufio.NewReader(r),
		lastDate: nil,
		order:    order,
	}
}

func (p *Parser) All() ([]Entry, error) {
	result := make([]Entry, 0)
	for {
		event, err := p.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		} else {
			result = append(result, *event)
		}
	}

	return result, nil
}

func (p *Parser) Next() (*Entry, error) {
	line, prefix, err := p.reader.ReadLine()
	if err != nil {
		return nil, err
	} else if prefix {
		return nil, fmt.Errorf("line too long")
	} else if len(line) == 0 || line[0] == '!' {
		return p.Next()
	}

	v := new(eventTarget)
	if err := eventParser.ParseBytes("***", line, v); err != nil {
		return nil, fmt.Errorf("event: %v", err)
	}

	dt := timeutil.Datetime{}
	if v.Date == nil && p.lastDate == nil {
		return nil, fmt.Errorf("no date context")
	} else if v.Date != nil {
		date, err := timeutil.ParseDate(v.Date.Left, v.Date.Middle, v.Date.Right, p.order)
		date.Weekday, _ = timeutil.ParseWeekday(v.Date.Weekday)
		if err != nil {
			return nil, fmt.Errorf("date: %w", err)
		}
		dt.Date = *date
		p.lastDate = date
	} else {
		dt.Date = *p.lastDate
	}

	if v.Body == nil {
		return p.Next()
	}

	if v.Time != nil {
		time, err := timeutil.ParseTime(v.Time.Hour, v.Time.Minute)
		if err != nil {
			return nil, fmt.Errorf("time: %w", err)
		}
		dt.Time = *time
	}

	dur := time.Duration(0)
	if v.Duration != nil {
		dur, err = time.ParseDuration(*v.Duration)
		if err != nil {
			return nil, fmt.Errorf("duration: %w", err)
		}
	}

	return &Entry{
		Timestamp:   dt,
		Duration:    dur,
		Description: *v.Body,
	}, nil

}

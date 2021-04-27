package diary

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/leotaku/maiden/timeutil"
)

func FromVEvent(ve *ics.VEvent) (*Entry, error) {
	summary := ve.GetProperty(ics.ComponentPropertySummary)
	if summary == nil || summary.Value == "" {
		return nil, fmt.Errorf("missing event summary")
	}

	rrule := ve.GetProperty(ics.ComponentPropertyRrule)
	if rrule != nil && rrule.Value != "" {
		return nil, fmt.Errorf("recurring events are not currently supported")
	}

	var datetime timeutil.Datetime
	start, err := getDatetimeProperty(ve, ics.ComponentPropertyDtStart)
	if err != nil {
		start, err = getDateProperty(ve, ics.ComponentPropertyDtStart)
		if err != nil {
			return nil, fmt.Errorf("start: %w", err)
		} else {
			datetime = timeutil.FromDate(*start)
		}
	} else {
		datetime = timeutil.FromAbsolute(*start)
	}

	end, err := getDatetimeProperty(ve, ics.ComponentPropertyDtEnd)
	if err != nil {
		end, err = getDateProperty(ve, ics.ComponentPropertyDtEnd)
		if err != nil {
			end = start
		}
	}

	return &Entry{
		Datetime:    datetime,
		Duration:    end.Sub(*start),
		Description: summary.Value,
	}, nil
}

func (e Entry) ToVEvent(loc *time.Location) *ics.VEvent {
	ve := new(ics.VEvent)
	ve.SetSummary(e.Description)
	ve.SetCreatedTime(time.Now())
	ve.SetDtStampTime(time.Now())
	ve.SetModifiedAt(time.Now())
	ve.SetLocation("")
	ve.SetTimeTransparency("Opaque")
	ve.SetSequence(0)
	ve.SetStatus(ics.ObjectStatusConfirmed)
	ve.SetDescription("")

	// Start and End
	start := e.Start(loc)
	end := e.End(loc)
	if e.Datetime.Time.Hour != nil && e.Datetime.Time.Minute != nil {
		setDatetimeProperty(ve, ics.ComponentPropertyDtStart, start)
		setDatetimeProperty(ve, ics.ComponentPropertyDtEnd, end)
	} else if e.Duration%(time.Hour*24) != 0 {
		if e.Duration == 0 {
			end = end.Add(time.Hour * 24)
		}
		setDatetimeProperty(ve, ics.ComponentPropertyDtStart, start)
		setDatetimeProperty(ve, ics.ComponentPropertyDtEnd, end)
	} else {
		setDateProperty(ve, ics.ComponentPropertyDtStart, start)
		setDateProperty(ve, ics.ComponentPropertyDtEnd, end)
	}

	// Default reminder
	va := ve.AddAlarm()
	va.SetAction("DISPLAY")
	va.SetProperty(ics.ComponentPropertyDescription, "This is an event reminder")
	va.SetTrigger("-P0DT0H30M0S")

	// Repeat
	rrule := e.RRule(loc)
	if rrule != "" {
		ve.AddRrule(rrule)
	}

	return ve
}

func getDatetimeProperty(ve *ics.VEvent, prop ics.ComponentProperty) (*time.Time, error) {
	p := ve.GetProperty(prop)
	if p != nil && p.Value != "" {
		loc := time.UTC
		params := p.ICalParameters["TZID"]
		if len(params) == 1 {
			loc, _ = time.LoadLocation(params[0])
		}

		time, err := time.ParseInLocation(timeutil.DateTimeFormat, p.Value, loc)
		if err != nil {
			return nil, fmt.Errorf("parse: %w", err)
		} else {
			return &time, nil
		}
	} else {
		return nil, fmt.Errorf("missing property: %v", prop)
	}
}

func getDateProperty(ve *ics.VEvent, prop ics.ComponentProperty) (*time.Time, error) {
	p := ve.GetProperty(prop)
	if p != nil && p.Value != "" {
		time, err := time.Parse(timeutil.DateFormat, p.Value)
		if err != nil {
			return nil, fmt.Errorf("parse: %w", err)
		}
		return &time, nil
	} else {
		return nil, fmt.Errorf("missing property: %v", prop)
	}
}

func setDatetimeProperty(ve *ics.VEvent, prop ics.ComponentProperty, d time.Time) {
	ve.SetProperty(prop, d.UTC().Format(timeutil.DateTimeFormat))
}

func setDateProperty(ve *ics.VEvent, prop ics.ComponentProperty, d time.Time) {
	ve.SetProperty(prop, d.Format(timeutil.DateFormat), &ics.KeyValues{
		Key:   "VALUE",
		Value: []string{"DATE"},
	})
}

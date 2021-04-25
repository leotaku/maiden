package diary

import (
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/leotaku/maiden/timeutil"
)

func (e Entry) ToICALEvent(loc *time.Location) *ics.VEvent {
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
		setLocalDateTime(ve, ics.ComponentPropertyDtStart, start)
		setLocalDateTime(ve, ics.ComponentPropertyDtEnd, end)
	} else {
		setLocalDate(ve, ics.ComponentPropertyDtStart, start)
		setLocalDate(ve, ics.ComponentPropertyDtEnd, end)
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

func setLocalDateTime(ve *ics.VEvent, prop ics.ComponentProperty, d time.Time) {
	ve.SetProperty(prop, d.UTC().Format(timeutil.DateTimeFormat))
}

func setLocalDate(ve *ics.VEvent, prop ics.ComponentProperty, d time.Time) {
	ve.SetProperty(prop, d.UTC().Format(timeutil.DateFormat), &ics.KeyValues{
		Key:   "VALUE",
		Value: []string{"DATE"},
	})
}

package caldav

import (
	"encoding/xml"
	"fmt"
	"strings"

	ics "github.com/arran4/golang-ical"
)

const getCtagXML = `
<d:propfind xmlns:d="DAV:" xmlns:cs="http://calendarserver.org/ns/">
  <d:prop>
     <cs:getctag />
  </d:prop>
</d:propfind>`

const getTimezoneXML = `
<d:propfind xmlns:d="DAV:" xmlns:c="urn:ietf:params:xml:ns:caldav">
  <d:prop>
     <c:calendar-timezone />
  </d:prop>
</d:propfind>`

const getEventsXML = `
<c:calendar-query xmlns:d="DAV:" xmlns:c="urn:ietf:params:xml:ns:caldav">
    <d:prop>
        <d:getetag />
        <c:calendar-data />
    </d:prop>
    <c:filter>
        <c:comp-filter name="VCALENDAR" />
    </c:filter>
</c:calendar-query>`

type multistatus struct {
	XMLName   xml.Name   `xml:"multistatus"`
	Responses []response `xml:"response"`
}

type response struct {
	Href   string `xml:"href"`
	Status string `xml:"propstat>status"`
	Props  props  `xml:"propstat>prop"`
}

type props struct {
	GetCtag          string `xml:"getctag"`
	GetEtag          string `xml:"getetag"`
	CalendarData     string `xml:"calendar-data"`
	CalendarTimezone string `xml:"calendar-timezone"`
}

func (rsp response) toEvent() (*ics.VEvent, error) {
	r := strings.NewReader(rsp.Props.CalendarData)
	cal, err := ics.ParseCalendar(r)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	} else if len(cal.Events()) != 1 {
		return nil, fmt.Errorf("not a singleton")
	}

	return cal.Events()[0], nil
}

func (rsp response) toTimezone() (*ics.VTimezone, error) {
	r := strings.NewReader(rsp.Props.CalendarTimezone)
	cal, err := ics.ParseCalendar(r)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	} else if len(cal.Components) != 1 {
		return nil, fmt.Errorf("not a singleton")
	}

	if vt, ok := cal.Components[0].(*ics.VTimezone); ok {
		return vt, nil
	} else {
		return nil, fmt.Errorf("not a timezone singleton")
	}
}

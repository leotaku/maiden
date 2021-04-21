package caldav

import (
	"fmt"

	ics "github.com/arran4/golang-ical"
)

type Client struct {
	eventIds idmap
	events   []*ics.VEvent
	builder
}

func (c *Client) ProviderName() string {
	return c.builder.providerName
}

func (c *Client) Events() []*ics.VEvent {
	return c.events
}

func (c *Client) Refetch() error {
	v, err := c.multistatus("REPORT", c.calendarPath, getEventsXML, 1)
	if err != nil {
		return fmt.Errorf("multistatus: %w", err)
	}

	result := make([]*ics.VEvent, 0)
	ids := make(idmap)
	for _, rsp := range v.Responses {
		ve, err := rsp.toEvent()
		if err != nil {
			return fmt.Errorf("event: %w", err)
		}
		result = append(result, ve)
		ids[ve.Id()] = idval{
			Href: rsp.Href,
			Etag: rsp.Props.GetEtag,
		}
	}

	c.events = result
	c.eventIds = ids
	return nil
}

func (c *Client) Put(ve *ics.VEvent) error {
	s := c.FinalizeEvent(ve).Serialize()
	href := c.eventHref(ve)
	req, err := c.prepareRequest("PUT", href, s)
	if err != nil {
		return fmt.Errorf("prepare: %v", err)
	}

	req.Header["Content-Type"] = []string{"text/calendar; charset=utf-8"}
	if etag := c.eventETag(ve); etag != "" {
		req.Header["If-Match"] = []string{etag}
	} else {
		req.Header["If-None-Match"] = []string{"*"}
	}

	rsp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 201 {
		return fmt.Errorf("status: %v", rsp.Status)
	}

	etag, err := getEtag(c.http, rsp)
	if err != nil {
		return fmt.Errorf("etag: %v", rsp.Status)
	}

	c.eventIds[ve.Id()] = idval{
		Href: href,
		Etag: etag,
	}

	return nil
}

func (c *Client) Del(ve *ics.VEvent) error {
	req, err := c.prepareRequest("DELETE", c.eventHref(ve), "")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}

	rsp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	} else if rsp.StatusCode != 204 {
		return fmt.Errorf("status: %v", rsp.StatusCode)
	} else {
		delete(c.eventIds, ve.Id())
		return nil
	}

}

func (c *Client) FinalizeEvent(ve *ics.VEvent) *ics.Calendar {
	if ve.Id() == "" {
		ve.SetProperty(ics.ComponentPropertyUniqueId, c.genId())
	}

	cal := ics.NewCalendar()
	*cal.AddEvent(ve.Id()) = *ve
	start := ve.GetProperty(ics.ComponentPropertyDtStart)
	if start != nil {
		if id := start.ICalParameters["TZID"]; len(id) == 1 {
			cal.SetXWRTimezone(id[0])
		}
	}

	return cal
}

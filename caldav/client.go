package caldav

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Client struct {
	eventIds idmap
	events   []*ics.VEvent
	timezone *ics.VTimezone
	builder
}

func (c *Client) ProviderName() string {
	return c.builder.providerName
}

func (c *Client) Events() []*ics.VEvent {
	return c.events
}

func (c *Client) Timezone() (string, *time.Location) {
	return getVTimezoneInfo(c.timezone)
}

func (c *Client) Refetch() error {
	err := c.fetchDefaultTimezone()
	if err != nil {
		return err
	}

	return c.fetchEvents()
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
	if c.timezone != nil {
		cal.Components = append(cal.Components, c.timezone)
		name, loc := getVTimezoneInfo(c.timezone)
		if start := ve.GetProperty(ics.ComponentPropertyDtStart); start != nil {
			normalizeTimezone(start, name, loc)
		}
		if end := ve.GetProperty(ics.ComponentPropertyDtEnd); end != nil {
			normalizeTimezone(end, name, loc)
		}
	}

	return cal
}

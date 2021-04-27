package caldav

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/leotaku/maiden/timeutil"
)

type (
	idmap = map[idkey]idval
	idkey = string
)

type idval struct {
	Href string
	Etag string
}

func (c *Client) fetchDefaultTimezone() error {
	v, err := c.multistatus("PROPFIND", c.calendarPath, getTimezoneXML, 0)
	if err != nil {
		return fmt.Errorf("multistatus: %v", err)
	}

	if len(v.Responses) == 1 {
		vt, err := v.Responses[0].toTimezone()
		if err != nil {
			return fmt.Errorf("timezone: %w", err)
		}
		c.timezone = vt
	}

	return nil
}

func (c *Client) fetchEvents() error {
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

func (c *Client) multistatus(method, href, body string, depth int) (*multistatus, error) {
	req, err := c.prepareRequest(method, href, body)
	if err != nil {
		return nil, fmt.Errorf("prepare: %v", err)
	}

	req.Header["Depth"] = []string{strconv.Itoa(depth)}
	req.Header["Prefer"] = []string{"return-minimal"}
	req.Header["Content-Type"] = []string{"application/xml; charset=utf-8"}

	rsp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	} else if rsp.StatusCode != 207 {
		return nil, fmt.Errorf("status: %v", rsp.StatusCode)
	}
	defer rsp.Body.Close()

	v := new(multistatus)
	err = xml.NewDecoder(rsp.Body).Decode(v)
	return v, err
}

func (c *Client) prepareRequest(method, href, body string) (*http.Request, error) {
	url, err := c.hostURL.Parse(href)
	if err != nil {
		return nil, fmt.Errorf("url: %w", err)
	}
	r := strings.NewReader(body)

	return http.NewRequest(method, url.String(), r)
}

func (c *Client) eventHref(ve *ics.VEvent) string {
	if v, ok := c.eventIds[ve.Id()]; ok {
		return v.Href
	} else {
		return path.Join(c.calendarPath, ve.Id()) + ".ics"
	}
}

func (c *Client) eventETag(ve *ics.VEvent) string {
	if v, ok := c.eventIds[ve.Id()]; ok {
		return v.Etag
	} else {
		return ""
	}
}

func (c *Client) genId() string {
	b := make([]byte, 16)
	rand.Read(b)
	id := fmt.Sprintf("%.16x@%v", b, c.providerName)

	if _, ok := c.eventIds[id]; ok {
		return c.genId()
	} else {
		return id
	}
}

func getEtag(c *http.Client, rsp *http.Response) (string, error) {
	if etag := rsp.Header.Get("ETag"); etag != "" {
		return etag, nil
	}

	req, err := http.NewRequest("HEAD", rsp.Request.URL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("prepare: %w", err)
	}

	if rsp, err := c.Do(req); err != nil {
		return "", err
	} else if etag := rsp.Header.Get("ETag"); etag == "" {
		return "", fmt.Errorf("empty")
	} else {
		return etag, nil
	}
}

func normalizeTimezone(prop *ics.IANAProperty, name string, loc *time.Location) {
	ploc := time.UTC
	if tzid := prop.ICalParameters["TZID"]; len(tzid) == 1 {
		ploc, _ = time.LoadLocation(tzid[0])
	}
	if time, err := time.ParseInLocation(timeutil.DateTimeFormat, prop.Value, ploc); err == nil {
		prop.Value = time.In(loc).Format(timeutil.DateTimeFormat)
		prop.ICalParameters["TZID"] = []string{name}
	}
}

func getVTimezoneInfo(vt *ics.VTimezone) (string, *time.Location) {
	for _, prop := range vt.Properties {
		if prop.IANAToken == string(ics.PropertyTzid) {
			time, err := time.LoadLocation(prop.Value)
			if err == nil {
				return prop.Value, time
			}
		}
	}

	return "UTC", nil
}

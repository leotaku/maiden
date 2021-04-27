package caldav

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"strings"

	ics "github.com/arran4/golang-ical"
)

type (
	idmap = map[idkey]idval
	idkey = string
)

type idval struct {
	Href string
	Etag string
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


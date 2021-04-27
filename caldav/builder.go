package caldav

import (
	"net/http"
	"net/url"
)

type builder struct {
	http         *http.Client
	hostURL      *url.URL
	calendarPath string
	providerName string
}

func NewBuilder() *builder {
	return &builder{
		http:         http.DefaultClient,
		providerName: "maiden@example.com",
	}
}

func (b *builder) BuildAndInit() (*Client, error) {
	c := &Client{
		eventIds: make(idmap),
		builder:  *b,
	}

	if err := c.Update(); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

func (b *builder) WithHostURL(url *url.URL) *builder {
	b.hostURL = url
	return b
}

func (b *builder) WithCalendarPath(path string) *builder {
	b.calendarPath = path
	return b
}

func (b *builder) WithHttp(http *http.Client) *builder {
	b.http = http
	return b
}

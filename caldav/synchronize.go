package caldav

import (
	"fmt"
	"time"
)

func (c *Client) Sync(interval time.Duration) *Syncer {
	return &Syncer{
		client:   c,
		interval: interval,
		kill:     make(chan struct{}, 1),
	}
}

func (c *Client) GetCtag() (string, error) {
	v, err := c.multistatus("PROPFIND", c.calendarPath, getCtagXML, 0)
	if err != nil {
		return "", fmt.Errorf("multistatus: %w", err)
	} else if len(v.Responses) != 1 {
		return "", fmt.Errorf("not exactly one response")
	} else {
		return v.Responses[0].Props.GetCtag, nil
	}
}

type Syncer struct {
	client   *Client
	interval time.Duration
	kill     chan struct{}
}

func (s *Syncer) Wait(ctag string) <-chan string {
	s.kill <- struct{}{}
	close(s.kill)
	s.kill = make(chan struct{}, 1)

	update := make(chan string)
	go func() {
		for {
			ct, err := s.client.GetCtag()
			if err != nil || ct == ctag {
				select {
				case <-time.After(s.interval):
				case <-s.kill:
					close(update)
					return
				}
			} else {
				update <- ct
				close(update)
				return
			}
		}
	}()

	return update
}

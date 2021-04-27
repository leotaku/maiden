package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/leotaku/maiden/caldav"
	"github.com/leotaku/maiden/diary"
	"github.com/leotaku/maiden/oauth"
	"github.com/leotaku/maiden/timeutil"
)

func loadGoogleClient(calID, dataHome string) (*caldav.Client, error) {
	life := oauth.NewLifecycle(dataHome)
	config, tok, err := life.Load()
	if err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}

	href := fmt.Sprintf("/caldav/v2/%v/events", calID)
	http := config.Client(context.TODO(), tok)
	b := caldav.NewBuilder()
	b.WithHttp(http)
	b.WithHostURL(googleURL)
	b.WithCalendarPath(href)
	client, err := b.BuildAndInit()
	if err != nil {
		return nil, fmt.Errorf("init: %w", err)
	}

	return client, nil
}

func validateDateStyle(style string) (timeutil.DateOrder, error) {
	switch strings.ToUpper(style) {
	case "ISO":
		return diary.ISO, nil
	case "AMERICAN":
		return diary.American, nil
	case "EUROPEAN":
		return diary.European, nil
	default:
		return nil, fmt.Errorf("invalid date order: %v", style)
	}
}

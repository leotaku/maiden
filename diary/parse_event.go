package diary

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

const (
	shortWday = `Mon|Tue|Wed|Thu|Fri|Sat|Sun`
	longWday  = `Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday`
	anyWday   = longWday + `|` + shortWday
)

const (
	shortMonth = `Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec`
	longMonth  = `January|February|March|April|May|June|July|August|September|October|November|December`
	anyMonth   = longMonth + `|` + shortMonth
	datePart   = anyMonth + `|[*]|[0-9]+`
)

var def = stateful.MustSimple([]stateful.Rule{
	{Name: "TP", Pattern: "[0-9]+:[0-9]+"},
	{Name: "DUR", Pattern: "([0-9]+[smh])+"},
	{Name: "WD", Pattern: anyWday},
	{Name: "DP", Pattern: datePart},
	{Name: "separator", Pattern: "[^[:alnum:]*]+?"},
	{Name: "Other", Pattern: ".*"},
})

type eventTarget struct {
	Date     *dateTarget `parser:"@@?"`
	Time     *timeTarget `parser:"@TP?"`
	Duration *string     `parser:"@DUR?"`
	Body     *string     `parser:"@Other?"`
}

type dateTarget struct {
	Left    string `parser:"@DP?"`
	Middle  string `parser:"@DP?"`
	Right   string `parser:"@DP?"`
	Weekday string `parser:"@WD?"`
}

type timeTarget struct {
	Hour   string
	Minute string
}

func (t *timeTarget) Capture(values []string) error {
	s := strings.Join(values, "")
	split := strings.Split(s, ":")
	if len(split) != 2 {
		return fmt.Errorf("not a time: %v", s)
	} else {
		t.Hour = split[0]
		t.Minute = split[1]
	}

	return nil
}

var eventParser = participle.MustBuild(
	&eventTarget{},
	participle.Lexer(def),
)

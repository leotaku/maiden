package diary_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/leotaku/maiden/diary"
)

const exampleDiary = `
2021/04/10   Important meeting
2021-04-11 08:00 Office work
April 19.    My birthday
 10:30 1h    Bake cake
 12:00       Jump off a cliff
Monday 10:00 Bring out the trash`

func TestParseLine(t *testing.T) {
	r := strings.NewReader(exampleDiary)
	parser := diary.NewParser(r, diary.ISO)
	events, err := parser.All()
	for _, event := range events {
		ve := event.ToICALEvent(time.Local)
		fmt.Println(ve.Serialize())
	}

	if err != nil {
		t.Error(err)
	}
}

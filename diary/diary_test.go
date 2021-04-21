package diary_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/leotaku/maiden/diary"
)

const exampleDiary = `
2021/10/10   Important meeting
April 19.    My birthday
 10:30 1h    Bake cake
 12:00       Jump off a cliff
Monday 10:00 Bring out the trash`

func TestParseLine(t *testing.T) {
	r := strings.NewReader(exampleDiary)
	parser := diary.NewParser(r, diary.ISO)
	events, err := parser.All()
	fmt.Println(err)
	for _, event := range events {
		ve := event.ToICALEvent(time.Local)
		fmt.Println(ve.Serialize())
	}

	t.Fail()
}

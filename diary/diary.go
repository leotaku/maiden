package diary

import (
	"fmt"
	"io"
	"os"

	"github.com/leotaku/maiden/timeutil"
)

var (
	ISO      = timeutil.ISO
	American = timeutil.American
	European = timeutil.European
)

type Diary struct {
	filename string
	entries  []Entry
	order    timeutil.DateOrder
}

func NewDiary(filename string, order timeutil.DateOrder) (*Diary, error) {
	d := Diary{
		filename: filename,
		entries:  make([]Entry, 0),
		order:    order,
	}

	if err := d.Update(); err != nil {
		return nil, fmt.Errorf("init: %w", err)
	} else {
		return &d, nil
	}
}

func (d *Diary) Entries() []Entry {
	return d.entries
}

func (d *Diary) Add(e Entry) error {
	if err := writeDiaryEntry(e, d.filename, d.order); err != nil {
		return err
	}
	d.entries = append(d.entries, e)
	return nil
}

func (d *Diary) Update() error {
	r, err := os.Open(d.filename)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	entries, err := NewParser(r, d.order).All()
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	d.entries = entries

	return nil
}

func writeDiaryEntry(entry Entry, filename string, order timeutil.DateOrder) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	separator := ""
	if bs[len(bs)-1] != '\n' {
		separator = "\n"
	}
	if _, err := f.WriteString(separator + entry.Format(order)); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return f.Sync()
}

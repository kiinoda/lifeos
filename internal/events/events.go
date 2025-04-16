package events

import (
	"fmt"
	"time"
)

type DayOfWeek int

const (
	Mon DayOfWeek = iota
	Tue
	Wed
	Thu
	Fri
	Sat
	Sun
)

type Event struct {
	Days      [7]string
	Time      time.Time
	Signifier string
	Desc      string
}

func NewEvent(line []any) (Event, error) {
	e := Event{}
	for day := Mon; day <= Sun; day++ {
		if line[day] != nil {
			e.Days[day] = line[day].(string)
		}
	}
	if line[8] != nil {
		if t, err := time.Parse("1504", line[8].(string)); err == nil {
			now := time.Now()
			e.Time = time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				t.Hour(),
				t.Minute(),
				0, 0,
				time.Local,
			)
		}
	}
	if line[9] != nil {
		if d, ok := line[9].(string); ok {
			e.Desc = d
		} else {
			return e, fmt.Errorf("Could not extract event description for %v", line[9])
		}
	}
	return e, nil
}

func (e Event) GetTimePlaceholder() string {
	if e.Time.IsZero() {
		return "----"
	}
	return e.Time.Format("1504")
}

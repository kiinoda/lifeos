package events

import (
	"fmt"
	"strings"
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

type ScheduledEvent struct {
	Alertable bool
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

func NewScheduledEvent(line []any) (ScheduledEvent, error) {
	event := ScheduledEvent{}
	if silent, ok := line[0].(string); ok {
		if strings.Trim(silent, " ") == "" {
			event.Alertable = true
		}
	} else {
		return event, fmt.Errorf("Could not extract event status for %v", line[0])
	}

	for _, format := range []string{"200601", "20060102", "2006Jan", "2006Jan02", "2006Jan2"} {
		if t, err := time.Parse(format, line[1].(string)); err == nil {
			event.Time = t
			break
		}
	}

	if line[3] != nil {
		if d, ok := line[3].(string); ok {
			event.Desc = d
		} else {
			return event, fmt.Errorf("Could not extract event description for %v", line[3])
		}
	}

	return event, nil
}

func (e Event) GetTimePlaceholder() string {
	if e.Time.IsZero() {
		return "----"
	}
	return e.Time.Format("1504")
}

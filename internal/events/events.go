package events

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var alertTimePattern = regexp.MustCompile(`^(\d{4})\s*(.*)$`)

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
	Recurring bool
	Alert     bool
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
	isYearKnown := false

	event := ScheduledEvent{}
	if silent, ok := line[0].(string); ok {
		if strings.Trim(silent, " ") == "" {
			event.Alertable = true
		}
	} else {
		return event, fmt.Errorf("Could not extract event status for %v", line[0])
	}

	event.Recurring = false
	event.Alert = false
	if marker, ok := line[2].(string); ok {
		trimmed := strings.TrimSpace(marker)
		if strings.EqualFold(trimmed, "r") {
			event.Recurring = true
		} else if strings.EqualFold(trimmed, "a") {
			event.Alert = true
		}
	}

	for _, format := range []string{"0102", "200601", "20060102", "2006Jan", "2006Jan02", "2006Jan2"} {
		eventTime := line[1].(string)
		if len(eventTime) == 4 {
			eventTime = strconv.Itoa(time.Now().Year()) + eventTime
		} else {
			isYearKnown = true
		}
		if t, err := time.Parse(format, eventTime); err == nil {
			event.Time = t
			break
		}
	}

	if line[3] != nil {
		if d, ok := line[3].(string); ok {
			event.Desc = d
			if event.Recurring {
				if isYearKnown {
					event.Desc = "(" + strconv.Itoa(event.Time.Year()) + ") " + event.Desc
				} else {
					event.Desc = "(R) " + event.Desc
				}
			}
		} else {
			return event, fmt.Errorf("Could not extract event description for %v", line[3])
		}
	}

	if event.Alert {
		if m := alertTimePattern.FindStringSubmatch(event.Desc); m != nil {
			if t, err := time.Parse("1504", m[1]); err == nil {
				y, mo, d := event.Time.Date()
				event.Time = time.Date(y, mo, d, t.Hour(), t.Minute(), 0, 0, time.Local)
				event.Desc = m[2]
			}
		}
	}

	// For recurring events, normalize to current year first, then advance to next year if more than 2 days past
	if event.Recurring {
		_, m, d := event.Time.Date()
		currentYearTime := time.Date(time.Now().Year(), m, d, 0, 0, 0, 0, event.Time.Location())
		if currentYearTime.Before(time.Now().Add(-48 * time.Hour)) {
			event.Time = time.Date(time.Now().Year()+1, m, d, 0, 0, 0, 0, event.Time.Location())
		} else {
			event.Time = currentYearTime
		}
	}

	return event, nil
}

func (se ScheduledEvent) ToEvent(dayOfWeek time.Weekday) Event {
	e := Event{
		Time: se.Time,
		Desc: se.Desc,
	}
	e.Days[dayOfWeek] = "x"
	return e
}

func (e Event) GetTimePlaceholder() string {
	if e.Time.IsZero() {
		return "----"
	}
	return e.Time.Format("1504")
}

package utils

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	schedule *cron.SpecSchedule
}

var parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func ParseCron(spec string) (*Cron, error) {
	schedule, err := parser.Parse(spec)
	if err != nil {
		return nil, err
	}
	if specSchedule, ok := schedule.(*cron.SpecSchedule); !ok {
		return nil, fmt.Errorf("bad cron %s", spec)
	} else {
		return &Cron{schedule: specSchedule}, nil
	}
}

func (c *Cron) Next(t time.Time) time.Time {
	return c.schedule.Next(t)
}

// Prev returns the previous time this schedule is activated, less than the given
// time.  If no time can be found to satisfy the schedule, return the zero time.
func (c *Cron) Prev(t time.Time) time.Time {
	// General approach is based on approach to Next() implementation

	origLocation := t.Location()
	loc := c.schedule.Location
	if loc == time.Local {
		loc = t.Location()
	}
	if c.schedule.Location != time.Local {
		t = t.In(c.schedule.Location)
	}

	// Start at the previous second
	t = t.Add(-1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// If no time is found within five years, return zero.
	yearLimit := t.Year() - 5

WRAP:
	if t.Year() < yearLimit {
		return time.Time{}
	}

	// Find the first applicable month.
	// If it's this month, then do nothing.
	for 1<<uint(t.Month())&c.schedule.Month == 0 {
		// set t to the last second of the previous month
		t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		t = t.Add(-1 * time.Second)

		// Wrapped around.
		if t.Month() == time.December {
			goto WRAP
		}
	}

	// Now get a day in that month.
	//
	// NOTE: This causes issues for daylight savings regimes where midnight does
	// not exist.  For example: Sao Paulo has DST that transforms midnight on
	// 11/3 into 1am. Handle that by noticing when the Hour ends up != 0.
	for !dayMatches(c.schedule, t) {
		// set t to the last second of the previous day
		saveMonth := t.Month()
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)

		// Notice if the hour is no longer midnight due to DST.
		// Add an hour if it's 23, subtract an hour if it's 1.
		if t.Hour() != 0 {
			if t.Hour() > 12 {
				t = t.Add(time.Duration(24-t.Hour()) * time.Hour)
			} else {
				t = t.Add(time.Duration(-t.Hour()) * time.Hour)
			}
		}

		t = t.Add(-1 * time.Second)

		if saveMonth != t.Month() {
			goto WRAP
		}
	}

	for 1<<uint(t.Hour())&c.schedule.Hour == 0 {
		// set t to the last second of the previous hour
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		t = t.Add(-1 * time.Second)

		if t.Hour() == 23 {
			goto WRAP
		}
	}

	for 1<<uint(t.Minute())&c.schedule.Minute == 0 {
		// set t to the last second of the previous minute
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, loc)
		t = t.Add(-1 * time.Second)

		if t.Minute() == 59 {
			goto WRAP
		}
	}

	for 1<<uint(t.Second())&c.schedule.Second == 0 {
		// set t to the previous second
		t = t.Add(-1 * time.Second)

		if t.Second() == 59 {
			goto WRAP
		}
	}

	return t.In(origLocation)
}

const (
	// Set the top bit if a star was included in the expression.
	starBit = 1 << 63
)

// dayMatches returns true if the schedule's day-of-week and day-of-month
// restrictions are satisfied by the given time.
func dayMatches(s *cron.SpecSchedule, t time.Time) bool {
	var (
		domMatch bool = 1<<uint(t.Day())&s.Dom > 0
		dowMatch bool = 1<<uint(t.Weekday())&s.Dow > 0
	)
	if s.Dom&starBit > 0 || s.Dow&starBit > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}

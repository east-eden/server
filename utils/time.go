package utils

import "time"

////////////////////////////////////////////////////////
// weekday
func PrevWeekday(d time.Weekday) time.Weekday { return (d - 1) % 7 }
func NextWeekday(d time.Weekday) time.Weekday { return (d + 1) % 7 }

func IsInSameDay(t1 time.Time, t2 time.Time, hour int) bool {
	if hour < 0 || hour >= 24 {
		return false
	}

	d := t1.Sub(t2)
	if (d / (time.Hour * 24)) > 0 {
		return false
	}

	t1 = t1.Add(-time.Hour * time.Duration(hour))
	t2 = t2.Add(-time.Hour * time.Duration(hour))

	return t1.Day() == t2.Day()
}

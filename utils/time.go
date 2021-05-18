package utils

import "time"

////////////////////////////////////////////////////////
// weekday
func PrevWeekday(d time.Weekday) time.Weekday { return (d - 1) % 7 }
func NextWeekday(d time.Weekday) time.Weekday { return (d + 1) % 7 }

func IsInSameDay(t1 time.Time, t2 time.Time, hour int) bool {

	return true
}

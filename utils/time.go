package utils

import (
	"os"
	"time"
)

const (
	OneYear           = time.Hour * 24 * 365
	HundredYear       = OneYear * 100
	LayoutRFC2822     = "2006/01/02 15:04:05"
	LayoutDateRFC2822 = "2006/01/02"

	// DefaultTimeZone 默认业务时区，运营活动、每日重置都按它计算
	DefaultTimeZone = "Asia/Shanghai"
)

// globalLoc 全局业务时区，可通过 TZ 环境变量覆盖
var globalLoc = mustLoadLocation(DefaultTimeZone)

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		// 部分精简镜像没有 tzdata，退化为固定东八区
		return time.FixedZone(name, 8*60*60)
	}

	return loc
}

func init() {
	timezone, ok := os.LookupEnv("TZ")
	if !ok || timezone == "" {
		return
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return
	}

	globalLoc = loc
}

// Location 返回当前业务时区
func Location() *time.Location {
	return globalLoc
}

// NowInLocation 返回业务时区下的当前时间
func NowInLocation() time.Time {
	return time.Now().In(globalLoc)
}

// GetTimestampNow 返回当前时间戳（毫秒）
func GetTimestampNow() int64 {
	return NowInLocation().UnixMilli()
}

// EnsureLocation 把时间转换到业务时区
func EnsureLocation(t time.Time) time.Time {
	return t.In(globalLoc)
}

// FromTimestamp 毫秒时间戳转业务时区时间
func FromTimestamp(milli int64) time.Time {
	return time.UnixMilli(milli).In(globalLoc)
}

// TodayZeroLocation 返回业务时区下的今日零点
func TodayZeroLocation() time.Time {
	now := NowInLocation()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, globalLoc)
}

// DateZeroFromTime 返回目标时间所在自然日的零点
func DateZeroFromTime(target time.Time) time.Time {
	year, month, day := target.In(globalLoc).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, globalLoc)
}

// DateZeroFromUnixMilli 返回目标毫秒时间戳所在自然日的零点
func DateZeroFromUnixMilli(milli int64) time.Time {
	return DateZeroFromTime(FromTimestamp(milli))
}

// FirstDayOfMonth 返回指定年月的第一天零点
func FirstDayOfMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, globalLoc)
}

// LongLongAfter 返回一个足够远的未来时间，用于表示“永不过期”
func LongLongAfter() time.Time {
	return time.Date(2100, 1, 18, 23, 59, 59, 0, globalLoc)
}

func ParseDateTime(date string) (time.Time, error) {
	return time.ParseInLocation(time.DateTime, date, globalLoc)
}

func ParseDateTimeRFC2822(date string) (time.Time, error) {
	return time.ParseInLocation(LayoutRFC2822, date, globalLoc)
}

func ParseDateRFC2822(date string) (time.Time, error) {
	return time.ParseInLocation(LayoutDateRFC2822, date, globalLoc)
}

// //////////////////////////////////////////////////////
// weekday
func PrevWeekday(d time.Weekday) time.Weekday { return (d - 1) % 7 }
func NextWeekday(d time.Weekday) time.Weekday { return (d + 1) % 7 }

// WeekDaySeven 返回 1-7 的星期表示（周一为 1，周日为 7）
func WeekDaySeven() int32 {
	weekday := NowInLocation().Weekday()
	if weekday == 0 {
		weekday = 7
	}

	return int32(weekday)
}

// IsInSameDay t1 must greater than t2
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

// IsInSameWeek t1 must greater than t2
func IsInSameWeek(t1 time.Time, t2 time.Time, hour int) bool {
	if hour < 0 || hour >= 24 {
		return false
	}

	d := t1.Sub(t2)
	if (d / (time.Hour * 24) / 7) > 0 {
		return false
	}

	// sub hour
	tt1 := t1.Add(-time.Hour * time.Duration(hour))
	tt2 := t2.Add(-time.Hour * time.Duration(hour))

	// sub weekday
	getSubWeekday := func(weekday time.Weekday) int {
		if weekday == time.Sunday {
			weekday += 7
		}

		return int(weekday - 1)
	}

	tt1 = tt1.Add(-time.Hour * 24 * time.Duration(getSubWeekday(tt1.Weekday())))
	tt2 = tt2.Add(-time.Hour * 24 * time.Duration(getSubWeekday(tt2.Weekday())))

	return tt1.Day() == tt2.Day()
}

package activity

import (
	"testing"
	"time"
)

// TestGetActivityTimeDatetime 固定时间模式
func TestGetActivityTimeDatetime(t *testing.T) {
	conf := &ActivityConf{
		StartMode: "datetime",
		Start:     "2024-01-01 10:00:00",
		EndMode:   "datetime",
		End:       "2024-01-02 10:00:00",
	}

	start, end := GetActivityTime(conf, time.Now(), time.Now(), 1)

	if start.Format(time.DateTime) != "2024-01-01 10:00:00" {
		t.Fatalf("start = %v", start)
	}

	if end.Format(time.DateTime) != "2024-01-02 10:00:00" {
		t.Fatalf("end = %v", end)
	}
}

// TestGetActivityTimeDuration 持续时间模式
func TestGetActivityTimeDuration(t *testing.T) {
	conf := &ActivityConf{
		StartMode: "datetime",
		Start:     "2024-01-01 10:00:00",
		EndMode:   "duration",
		End:       "24h",
	}

	start, end := GetActivityTime(conf, time.Now(), time.Now(), 1)
	if end.Sub(start) != 24*time.Hour {
		t.Fatalf("duration = %v, want 24h", end.Sub(start))
	}
}

// TestGetActivityTimeCron cron 模式取当前所处区间
func TestGetActivityTimeCron(t *testing.T) {
	conf := &ActivityConf{
		StartMode: "cron",
		Start:     "0 0 10 * * *", // 每天 10:00
		EndMode:   "duration",
		End:       "1h",
	}

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local)
	start, end := GetActivityTime(conf, now, now, 1)

	// 12:00 已过当天 10:00 的区间（10:00-11:00），应取下一次 10:00
	want := time.Date(2024, 1, 2, 10, 0, 0, 0, time.Local)
	if !start.Equal(want) {
		t.Fatalf("start = %v, want %v", start, want)
	}

	if end.Sub(start) != time.Hour {
		t.Fatalf("duration = %v, want 1h", end.Sub(start))
	}
}

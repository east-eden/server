package time2

import (
	"time"
)

var gTimecop *TimeCop = nil

// Deprecated: use Unix
func GetNowTS() int64 {
	return Unix()
}

// current timestamp
// time consume: time.Now().Unix() > atomic.LoadInt64(&_atomic_timestamp) > gTimecop.Unix() > gTimecop.ts
func Unix() int64 {
	return gTimecop.Unix()
}

func Now() time.Time {
	return gTimecop.Now()
}

func Freeze(t time.Time) {
	gTimecop.Freeze(t)
}

func Travel(t time.Time) {
	gTimecop.Travel(t)
}

func Scale(scale float64) {
	gTimecop.Scale(scale)
}

func Since(t time.Time) time.Duration {
	return gTimecop.Since(t)
}

func Sleep(d time.Duration) {
	gTimecop.Sleep(d)
}

func After(d time.Duration) <-chan time.Time {
	return gTimecop.After(d)
}

func Tick(d time.Duration) <-chan time.Time {
	return gTimecop.Tick(d)
}

func Return() {
	gTimecop.Return()
}

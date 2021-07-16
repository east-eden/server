// Package timewheel implements a timewheel algorithm that is suitable for large numbers of timers.
// It is based on golang channel broadcast mechanism.

package time2

import (
	"sync"
	"time"
)

var (
	timerMap map[time.Duration]*TimingWheel
	mapLock  = &sync.Mutex{}
	accuracy = 20 // means max 1/20 deviation
)

func init() {
	timerMap = make(map[time.Duration]*TimingWheel)
}

// After like time.After
func TimingWheelAfter(t time.Duration) <-chan struct{} {
	mapLock.Lock()
	defer mapLock.Unlock()
	if v, ok := timerMap[t]; ok {
		return v.After(t)
	}
	v := NewTimewheel(t/time.Duration(accuracy), accuracy+1)
	timerMap[t] = v
	return v.After(t)
}

// SetAccuracy sets the accuracy for the timewheel.
// low accuracy usually have better performance.
func SetAccuracy(a int) {
	accuracy = a
}

type TimingWheel struct {
	sync.Mutex
	interval   time.Duration
	ticker     *time.Ticker
	quit       chan struct{}
	maxTimeout time.Duration
	cs         []chan struct{}
	pos        int
}

func NewTimewheel(interval time.Duration, buckets int) *TimingWheel {
	w := &TimingWheel{
		interval:   interval,
		quit:       make(chan struct{}),
		pos:        0,
		maxTimeout: time.Duration(interval * (time.Duration(buckets - 1))),
		cs:         make([]chan struct{}, buckets),
		ticker:     time.NewTicker(interval),
	}
	for i := range w.cs {
		w.cs[i] = make(chan struct{})
	}
	go w.run()
	return w
}

func (w *TimingWheel) Stop() {
	close(w.quit)
}

// 误差在一个interval内
// timeline : ---w.pos-1<--{x}-->call After()<--{y}-->w.pos-----
// x + y == interval, y 即是误差
func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout > w.maxTimeout {
		timeout = w.maxTimeout
	} else if timeout < time.Second {
		timeout = time.Second
	}

	w.Lock()
	index := (w.pos + int(timeout/w.interval)) % len(w.cs)
	b := w.cs[index]
	w.Unlock()
	return b
}

func (w *TimingWheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) onTicker() {
	w.Lock()
	lastC := w.cs[w.pos]
	w.cs[w.pos] = make(chan struct{})
	w.pos = (w.pos + 1) % len(w.cs)
	w.Unlock()
	close(lastC)
}

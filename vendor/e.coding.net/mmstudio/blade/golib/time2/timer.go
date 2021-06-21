package time2

import (
	"fmt"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/golib/sync2"

	"e.coding.net/mmstudio/blade/golib/cron"
	"github.com/rs/zerolog/log"
)

const (
	DefaultTimerDomain = "timer-domain-system"
)

type TimerDispatcher interface {
	AfterFunc(d time.Duration, cb func()) *Timer
	CronFunc(cronExpr *cron.Expression, cb func()) *Cron

	RemoveAllTimer()
	RemoveAllTimerInDomain(domain string)
	AfterFuncInDomain(td time.Duration, cb func(), domain string) *Timer
}

// one dispatcher per goroutine (goroutine not safe)
type Dispatcher struct {
	ChanTimer     chan *Timer // one receiver, N senders
	runningTimers sync.Map
	queueLen      int
	stopChan      chan struct{}
	closeFlag     sync2.AtomicInt32
}

func NewDispatcher(l int) *Dispatcher {
	d := new(Dispatcher)
	d.queueLen = l
	d.start()
	return d
}

// for reopen
func (d *Dispatcher) start() {
	d.stopChan = make(chan struct{})
	d.ChanTimer = make(chan *Timer, d.queueLen)
	d.closeFlag.Set(0)
}
func (d *Dispatcher) Close() {
	if d.closeFlag.CompareAndSwap(0, 1) {
		close(d.stopChan)
		close(d.ChanTimer)

		//clear all timers
		for t := range d.ChanTimer {
			t.Cb()
		}
		d.RemoveAllTimer()
	}
}

func (d *Dispatcher) RemoveTimer(t *Timer) {
	d.runningTimers.Delete(t)
}
func (d *Dispatcher) RemoveAllTimer() {
	d.runningTimers.Range(func(key, value interface{}) bool {
		t := key.(*Timer)
		t.Stop()
		d.RemoveTimer(t)
		return true
	})
}
func (d *Dispatcher) RemoveAllTimerInDomain(domain string) {
	d.runningTimers.Range(func(key, value interface{}) bool {
		t := key.(*Timer)
		if t.domain != domain {
			return true
		}
		t.Stop()
		d.RemoveTimer(t)
		return true
	})
}

func (d *Dispatcher) AfterFunc(td time.Duration, cb func()) *Timer {
	return d.AfterFuncInDomain(td, cb, DefaultTimerDomain)
}
func (d *Dispatcher) AfterFuncInDomain(td time.Duration, cb func(), domain string) *Timer {
	t := new(Timer)
	t.cb = cb
	t.domain = domain
	t.t = time.AfterFunc(td, func() {
		// callback from another goroutine
		select {
		// FIRST read from no buffer chan, even closed, will return false
		case <-d.stopChan:
			return
		default:
			d.ChanTimer <- t
		}
	})
	d.runningTimers.Store(t, struct{}{})
	log.Debug().Msg(fmt.Sprintf("Timer Dispatcher add AfterFuncInDomain:%s after:%s", domain, td))
	return t
}

func (d *Dispatcher) CronFunc(cronExpr *cron.Expression, _cb func()) *Cron {
	c := new(Cron)

	now := time.Now()
	nextTime := cronExpr.Next(now)
	if nextTime.IsZero() {
		return c
	}

	// callback
	var cb func()
	cb = func() {
		defer _cb()
		now := time.Now()
		nextTime := cronExpr.Next(now)
		if nextTime.IsZero() {
			return
		}
		c.t = d.AfterFunc(nextTime.Sub(now), cb)
	}

	c.t = d.AfterFunc(nextTime.Sub(now), cb)
	return c
}

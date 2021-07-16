package time2

import "time"

// Timer
type Timer struct {
	t      *time.Timer
	domain string
	cb     func()
}

func (t *Timer) Stop() {
	t.t.Stop()
	t.cb = nil
}
func (t *Timer) Reset(d time.Duration) bool { return t.t.Reset(d) }

// call from Dispatcher's goroutine
func (t *Timer) Cb() {
	defer func() {
		t.cb = nil
		// todo should we catch panic?
	}()
	if t.cb != nil {
		t.cb()
	}
}

// Cron
type Cron struct {
	t *Timer
}

func (c *Cron) Stop() {
	if c.t != nil {
		c.t.Stop()
	}
}

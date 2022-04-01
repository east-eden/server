package ttlwait

import "time"

type TTLWait struct {
	t *time.Timer
	c chan bool
}

func NewTTLWait(d time.Duration) *TTLWait {
	return &TTLWait{
		t: time.NewTimer(d),
		c: make(chan bool, 1),
	}
}

func (tc *TTLWait) Done() {
	tc.c <- true
}

func (tc *TTLWait) Wait() {
	select {
	case <-tc.c:
		break
	case <-tc.t.C:
		break
	}

	tc.t.Stop()
}

package ttlwait

import (
	"sync/atomic"
	"time"
)

type TTLWaitGroup struct {
	t *time.Timer
	c chan bool
	n int32
}

func NewTTLWaitGroup(d time.Duration, n int32) *TTLWaitGroup {
	return &TTLWaitGroup{
		t: time.NewTimer(d),
		c: make(chan bool, 1),
		n: n,
	}
}

func (g *TTLWaitGroup) Add(n int32) {
	atomic.AddInt32(&g.n, n)
}

func (g *TTLWaitGroup) Done() {
	g.c <- true
}

func (g *TTLWaitGroup) Wait() (timeout bool) {
	defer g.t.Stop()

	for {
		select {
		case <-g.c:
			nn := atomic.AddInt32(&g.n, -1)
			if nn <= 0 {
				return
			}
		case <-g.t.C:
			timeout = true
			return
		}
	}
}

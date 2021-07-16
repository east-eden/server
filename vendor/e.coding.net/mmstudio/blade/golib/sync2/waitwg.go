package sync2

import (
	"sync"
	"time"
)

// WaitGroupTimeout is a wrapper of sync.WaitGroup that
// supports wait with timeout.
type WaitGroupTimeout struct {
	sync.WaitGroup
}

func (w *WaitGroupTimeout) WaitTimeout(timeout time.Duration) bool {
	return WaitTimeout(&w.WaitGroup, timeout)
}

// WaitTimeout performs a timed wait on a given sync.WaitGroup
// FIXME if timeout triggered, there will be goroutine leak.
func WaitTimeout(waitGroup *sync.WaitGroup, timeout time.Duration) bool {
	success := make(chan struct{})
	go func() {
		defer func() {
			// swallow any panics, as they'll just be from the channel
			// close if the timeout elapsed
			recover()
		}()
		defer close(success)
		waitGroup.Wait()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-success: // completed normally
		return false
	case <-timer.C:
		return true
	}
}

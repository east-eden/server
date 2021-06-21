package time2

import (
	"context"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/golib/suturev4"
	"e.coding.net/mmstudio/blade/golib/sync2"
)

// mock like ruby timecop
// TODO WH mock should not usable in prod

type TimeCop struct {
	ts        sync2.AtomicInt64
	running   sync2.AtomicBool
	closeChan chan struct{}
	closeFlag sync2.AtomicInt32
	stoken    suturev4.ServiceToken

	// for mock
	freezeTime time.Time
	travelTime time.Time
	scale      float64
	rw         *sync.RWMutex
	frozen     bool
	traveled   bool
}

// run as global will replace the global timecop and auto serve background
func NewTimeCop(asGlobal bool) *TimeCop {
	supervisorOnce.Do(func() {
		tmp := suturev4.New("time2", defaultSupervisorSpec)
		tmp.ServeBackground(context.Background())
		SetSkeletonSupervisor(tmp)
	})
	tc := &TimeCop{}
	tc.rw = new(sync.RWMutex)
	tc.closeChan = make(chan struct{})
	tc.stoken = suturev4.ServiceToken{}
	if asGlobal {
		tc.AsGlobal()
	}
	tc.Return()
	return tc
}

func (tc *TimeCop) AsGlobal() {
	if gTimecop != nil && gTimecop != tc {
		gTimecop.stopServeBackground()
	}
	gTimecop = tc
	tc.serveBackground()
}

func (tc *TimeCop) serveBackground() {
	if tc.stoken != (suturev4.ServiceToken{}) {
		return
	}
	tc.stoken = time2Supervisor.Add(tc)
}

func (tc *TimeCop) stopServeBackground() {
	if tc.stoken == (suturev4.ServiceToken{}) {
		return
	}
	_ = time2Supervisor.Remove(tc.stoken)
	tc.stoken = suturev4.ServiceToken{}
}

func (tc *TimeCop) Serve(ctx context.Context) error {
	tc.running.Set(true)
	tc.ts.Set(time.Now().Unix())
	c := time.Tick(1 * time.Second)
	for {
		select {
		case now := <-c:
			tc.ts.Set(now.Unix())
		case <-tc.closeChan:
			tc.running.Set(false)
			return suturev4.ErrDoNotRestart
		case <-ctx.Done():
			return suturev4.ErrDoNotRestart
		}
	}
}

func (tc *TimeCop) Stop() {
	if tc.closeFlag.CompareAndSwap(0, 1) {
		close(tc.closeChan)
	}
}

func (tc *TimeCop) Mocked() bool {
	return tc.frozen || tc.traveled
}

// warning: 最小单位是秒, time.Now()最小单位是纳秒
func (tc *TimeCop) Now() time.Time {
	if tc.Mocked() {
		tc.rw.RLock()
		defer tc.rw.RUnlock()
		if tc.frozen {
			return tc.freezeTime
		}
		if tc.traveled {
			return tc.freezeTime.Add(time.Duration(float64(time.Unix(tc.UnixWithoutMock(), 0).Sub(tc.travelTime)) * tc.scale))
		}
	}
	return time.Unix(tc.UnixWithoutMock(), 0)
}

func (tc *TimeCop) UnixWithoutMock() int64 {
	if tc.running.Get() {
		return tc.ts.Get()
	}
	return time.Now().Unix()
}

func (tc *TimeCop) Unix() int64 {
	if tc.Mocked() {
		return tc.Now().Unix()
	}
	return tc.UnixWithoutMock()
}

func (tc *TimeCop) Freeze(t time.Time) {
	tc.rw.Lock()
	defer tc.rw.Unlock()
	tc.freezeTime = t
	tc.frozen = true
}

func (tc *TimeCop) Scale(scale float64) {
	tc.scale = scale
	if !tc.traveled {
		tc.Travel(tc.Now())
	}
}

func (tc *TimeCop) Travel(t time.Time) {
	tc.rw.Lock()
	defer tc.rw.Unlock()
	tc.freezeTime = t
	tc.travelTime = time.Unix(tc.UnixWithoutMock(), 0)
	tc.traveled = true
}

func (tc *TimeCop) Since(t time.Time) time.Duration {
	return tc.Now().Sub(t)
}

func (tc *TimeCop) Sleep(d time.Duration) {
	if tc.Mocked() && d > 0 {
		tc.Travel(tc.Now().Add(d))
	} else {
		time.Sleep(d)
	}
}

func (tc *TimeCop) After(d time.Duration) <-chan time.Time {
	if tc.Mocked() {
		tc.Travel(tc.Now().Add(d))
		c := make(chan time.Time, 1)
		c <- tc.Now()
		return c
	} else {
		return time.After(d)
	}
}

func (tc *TimeCop) Tick(d time.Duration) <-chan time.Time {
	if tc.Mocked() {
		c := make(chan time.Time, 1)
		go func() {
			for {
				tc.Travel(tc.Now().Add(d))
				c <- tc.Now()
			}
		}()
		return c
	} else {
		return time.Tick(d)
	}
}

func (tc *TimeCop) Return() {
	tc.rw.Lock()
	defer tc.rw.Unlock()
	tc.frozen = false
	tc.traveled = false
	tc.scale = 1
}

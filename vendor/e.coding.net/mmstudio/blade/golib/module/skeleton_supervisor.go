package module

import (
	"context"
	"fmt"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/golib/suturev4"
)

const skeletonSupervisorTag = "skeleton"

var skeletonSupervisorOnce sync.Once
var defaultSupervisorSpec = suturev4.Spec{
	FailureDecay:     30,                              // 30 seconds
	FailureThreshold: 5,                               //5 failures
	FailureBackoff:   time.Duration(15) * time.Second, //15 seconds
	Timeout:          time.Duration(10) * time.Second, // 10 seconds
}

var skeletonSupervisor *suturev4.Supervisor

func SetSkeletonSupervisor(supervisor *suturev4.Supervisor) { skeletonSupervisor = supervisor }

// SupervisorServer adds current skeleton to skeleton supervisor.
// A timeout value of 0 means to wait forever if the server can not start
func (s *skeleton) SupervisorServer(timeout ...time.Duration) error {
	skeletonSupervisorOnce.Do(func() {
		if skeletonSupervisor == nil {
			tmp := suturev4.New(skeletonSupervisorTag, defaultSupervisorSpec)
			tmp.ServeBackground(context.Background())
			SetSkeletonSupervisor(tmp)
		}
	})

	s.sutureTokenMutex.Lock()
	defer s.sutureTokenMutex.Unlock()

	if s.sutureToken != (suturev4.ServiceToken{}) {
		return fmt.Errorf("suture token not empty,got:%d", s.sutureToken)
	}
	s.sutureToken = skeletonSupervisor.Add(s)

	// sync to wait service started
	var timeoutC <-chan time.Time
	if len(timeout) > 0 && timeout[0] > 0 {
		timer := time.NewTimer(timeout[0])
		defer timer.Stop()
		timeoutC = timer.C
	}
	select {
	case <-timeoutC:
		return fmt.Errorf("failed to start the skeleton:%s after:%s", s._name, timeout)
	case <-s.startedChan:
		return nil
	}
}

// SupervisorStopNoWait will remove current skeleton from the Supervisor, this returns without waiting for the skeleton to stop.
func (s *skeleton) SupervisorStopNoWait() (err error) { return s.stop(0, false) }

// SupervisorStop will remove current skeleton from the Supervisor
// A timeout value of 0 means to wait forever.
//
// Do not call this as action like this,will block forever.
// var s := NewSkeleton("test")
// s.AsyncInvoke(func() {
//	  s.SupervisorStop(0)
// })
func (s *skeleton) SupervisorStop(timeout time.Duration) (err error) { return s.stop(timeout, true) }

func (s *skeleton) stop(timeout time.Duration, wait bool) (err error) {
	s.sutureTokenMutex.Lock()
	defer s.sutureTokenMutex.Unlock()

	// suture层顺序处理，如果service已经重启，当前将service移除，会正常触发stopInner信号量
	// 如果在重启之前将service移除，则stoppedChanInner信号量无法得到触发

	if s.sutureToken == (suturev4.ServiceToken{}) {
		return nil
	}
	if wait {
		err = skeletonSupervisor.RemoveAndWait(s.sutureToken, timeout)
	} else {
		// no wait
		err = skeletonSupervisor.Remove(s.sutureToken)
	}
	if err != nil {
		return err
	}

	s.sutureToken = suturev4.ServiceToken{}

	// 当前处于terminate重启阶段,主动关闭
	if s.terminated.Get() {
		s.onStop(timeout)
		return nil
	}

	if !wait {
		// no wait
		return nil
	}
	// sync to wait service started
	var timeoutC <-chan time.Time
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		timeoutC = timer.C
	}
	select {
	case <-timeoutC:
		return fmt.Errorf("failed to stop the skeleton:%s after:%s", s._name, timeout)
	case <-s.stoppedChanInner:
		return nil
	}
}

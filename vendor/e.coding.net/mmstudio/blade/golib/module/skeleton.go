package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"e.coding.net/mmstudio/blade/golib/gid"
	"e.coding.net/mmstudio/blade/golib/gls"
	"e.coding.net/mmstudio/blade/golib/liveness"
	"e.coding.net/mmstudio/blade/golib/paniccatcher"
	"e.coding.net/mmstudio/blade/golib/suturev4"
	"e.coding.net/mmstudio/blade/golib/sync2"
	"e.coding.net/mmstudio/blade/golib/time2"
	"github.com/rs/zerolog/log"
)

type Action = func() (interface{}, error)
type ActionWithContext = func(context.Context) (interface{}, error)
type ActionNoReturn = func()
type ActionNoReturnWithContext = func(context.Context)

var DefaultTerminateStopTimeout = time.Duration(20) * time.Second
var liveTimeWheel = time2.NewTimewheel(time.Second*time.Duration(5), 10)
var syncCallTimer *time2.TimingWheel

func init() { SetSkeletonInvokeTimerMaxBucket(defaultTimerMaxBucket) }

// SetSkeletonInvokeTimerMaxBucket mainly for debug, not goroutine safe not call this while running
func SetSkeletonInvokeTimerMaxBucket(bucket int) {
	syncCallTimer = time2.NewTimewheel(time.Second, bucket)
}

type Skeleton interface {
	skeletonImpl

	IsRunning() bool

	SyncInvoke(fn Action) (interface{}, error)
	SyncInvokeWithContext(ctx context.Context, fn ActionWithContext) (interface{}, error)

	AsyncInvoke(fn ActionNoReturn) error
	AsyncInvokeWithContext(ctx context.Context, fn ActionNoReturnWithContext) error

	StoppedChan() chan struct{}
	TerminatedChan() chan struct{}
	SetEventHook(cc chan struct{}, hook func(Skeleton))

	// SupervisorServer 模式启动，此种模式下，如果skeleton异常重启会自动重新拉起
	SupervisorServer(timeout ...time.Duration) error
	SupervisorStop(timeout time.Duration) error
	SupervisorStopNoWait() error

	GID() unsafe.Pointer
	SkeletonID() uint64

	SetUserData(interface{})
	UserData() interface{}
	RestartCount() int

	time2.TimerDispatcher
}

type execItem struct {
	action  ActionWithContext
	ctx     context.Context
	timeout sync2.AtomicBool
	retChan chan *result
	from    unsafe.Pointer
}

type result struct {
	err    error
	result interface{}
}

var executeItemPool = &sync.Pool{
	New: func() interface{} {
		return &execItem{retChan: make(chan *result, 1)}
	},
}

func allocExecuteItem(ctx context.Context, action ActionWithContext) *execItem {
	ei := executeItemPool.Get().(*execItem)
	ei.ctx = ctx
	ei.action = action
	return ei
}

func freeExecuteItem(item *execItem) {
	item.action = nil // avoid lack of captured vars
	item.ctx = nil
	item.from = gid.Invalid
	item.timeout.Set(false)
	executeItemPool.Put(item)
}

func (s *skeleton) exec(ei *execItem) {
	if ei == nil {
		return
	}
	if ei.timeout.Get() {
		ei.retChan <- &result{result: nil, err: newSkeletonError(s).SetTimeout()}
		return
	}
	if err := ei.ctx.Err(); err != nil {
		if err == context.DeadlineExceeded {
			ei.retChan <- &result{result: nil, err: newSkeletonError(s).SetTimeout()}
		} else {
			ei.retChan <- &result{result: nil, err: err}
		}
		return
	}

	if !s.IsRunning() {
		ei.retChan <- &result{result: nil, err: newSkeletonError(s)}
		return
	}

	s.setActionFrom(ei.from)
	// fixme should we defer here? if the caller goroutine exist,the captured vars(pointer) could be changed.
	if s.spec.CatchPanic {

		// todo panic的item是否已经塞回队列等重启后再次执行？
		paniccatcher.Do(func() {
			ret, err := ei.action(ei.ctx)
			// should not be blocked, could use future'server goroutine do this,especially in async model
			ei.retChan <- &result{result: ret, err: err}
		}, func(p *paniccatcher.Panic) {
			event := log.Error()
			event = event.Str("skeleton", s.name())
			event = event.Interface("reason", p.Reason)
			if s.spec.PanicWithStack {
				event = event.Str("stack", string(debug.Stack()))
			}
			event.Msg("skeleton exec got panic")
			ei.retChan <- &result{result: nil, err: newSkeletonError(s).SetPanic(fmt.Sprintf("%s", p.Reason))}
		})
	} else {
		ret, err := ei.action(ei.ctx)
		ei.retChan <- &result{result: ret, err: err}
	}
	s.setActionFrom(gid.Invalid)
}

func (s *skeleton) state() string {
	state := &struct {
		Name        string
		ExecLen     int
		RingExecLen int
	}{
		Name:        s.name(),
		ExecLen:     len(s.execChan),
		RingExecLen: len(s.ringExecChan),
	}
	b, _ := json.Marshal(state)
	return string(b)
}

func (s *skeleton) selectAppendExecItem(ei *execItem, timeout time.Duration) error {
	select {
	case s.execChan <- ei:
	case <-syncCallTimer.After(timeout):
		// should not free the ei here, but if the caller goroutine exit, the ei may escaped from our pool
		ei.timeout.Set(true)
		return newSkeletonError(s).SetBlocked().SetTimeout()
	}
	return nil
}

// fixme when the skeleton closed normally, should return error immediately instead of waiting for timeout
func (s *skeleton) call(ctx context.Context, action ActionWithContext) (interface{}, error) {
	g := gid.Get()

	var callerSkeletonId unsafe.Pointer
	var caller Skeleton

	// 获取调用者的skeleton信息，可能并不存在
	if c, ok := gid2SkeletonMap.Load(g); ok {
		caller = c.(Skeleton)
		callerSkeletonId = caller.GID()
	}
	// 超时逻辑，如果context有提供的Deadline，获取最先到达的时间
	invokeTimeout := s.spec.InvokeTimeout
	if deadline, ok := ctx.Deadline(); ok {
		// ctx若有超时，一定在s.spec.InvokeTimeout之前被触发
		if duration := deadline.Sub(time.Now()); duration >= invokeTimeout {
			invokeTimeout += duration
		}
	}
	//fmt.Println(fmt.Sprintf("call from %p",g))
	//fmt.Println(fmt.Sprintf("call on %p",s.gid))
	//fmt.Println(fmt.Sprintf("caller %p",callerSkeletonId))
	selfGid := s.gid.Get()
	// 调用方来自于当前skeleton的协程，同步调用直接执行
	if g == selfGid && s.IsRunning() {
		// A 停止后启动，gid可能被另一个协程B获取，B调A会导致逻辑错误走到这里，
		return action(ctx)
	}

	// 被调用skeleton必须是执行状态，否则报错
	if !s.IsRunning() {
		return nil, newSkeletonError(s).SetClosed()
	}

	if g == callerSkeletonId {
		// 由另一个skeleton调动而来，可能存在环
		ei := allocExecuteItem(ctx, action)
		if s.spec.StrictOrderInvoke {
			// just try our best to detect ring call like A->B->C->A
			isRing := false
			if s.spec.DetectLoopInvoke {
				ringCall, stack := tryDetectRing(selfGid, callerSkeletonId, s.spec.StepDetectLoopInvoke)
				isRing = ringCall
				if isRing {
					s.spec.LoopInvokeTrigger(stack)
				}
			}
			if isRing {
				s.ringExecChan <- ei
				select {
				case ret := <-ei.retChan:
					freeExecuteItem(ei)
					return ret.result, ret.err
				case <-syncCallTimer.After(invokeTimeout):
					// should not free the ei here, but if the caller goroutine exit, the ei may escaped from our pool
					ei.timeout.Set(true)
					return nil, newSkeletonError(s).SetTimeout()
				case <-ctx.Done():
					if ctx.Err() == context.DeadlineExceeded {
						ei.timeout.Set(true)
						return nil, newSkeletonError(s).SetTimeout()
					} else {
						return nil, ctx.Err()
					}
				}
			} else {
				ei.from = callerSkeletonId
				select {
				case s.execChan <- ei:
				case <-syncCallTimer.After(invokeTimeout):
					ei.timeout.Set(true)
					return nil, newSkeletonError(s).SetTimeout().SetBlocked()
				case <-ctx.Done():
					if ctx.Err() == context.DeadlineExceeded {
						ei.timeout.Set(true)
						return nil, newSkeletonError(s).SetTimeout().SetBlocked()
					} else {
						return nil, ctx.Err()
					}
				}
				for {
					select {
					case ret := <-ei.retChan:
						freeExecuteItem(ei)
						return ret.result, ret.err
					case ei := <-caller.getRingExecChan(): // for ring：A->B->C->A
						caller.exec(ei)
					case <-syncCallTimer.After(invokeTimeout):
						ei.timeout.Set(true)
						return nil, newSkeletonError(s).SetTimeout()
					case <-ctx.Done():
						if ctx.Err() == context.DeadlineExceeded {
							ei.timeout.Set(true)
							return nil, newSkeletonError(s).SetTimeout()
						} else {
							return nil, ctx.Err()
						}
					}
				}
			}
		} else {
		executeLoop:
			for {
				select {
				case s.execChan <- ei:
					break executeLoop
				case ei := <-caller.getExecChan():
					caller.exec(ei)
				}
			}
			for {
				select {
				case ret := <-ei.retChan:
					freeExecuteItem(ei)
					return ret.result, ret.err
				case ei := <-caller.getExecChan(): // for ring：A->B->C->A
					caller.exec(ei)
				case <-syncCallTimer.After(invokeTimeout):
					ei.timeout.Set(true)
					return nil, newSkeletonError(s).SetTimeout()
				case <-ctx.Done():
					if ctx.Err() == context.DeadlineExceeded {
						ei.timeout.Set(true)
						return nil, newSkeletonError(s).SetTimeout()
					} else {
						return nil, ctx.Err()
					}
				}
			}
		}
	} else {
		ei := allocExecuteItem(ctx, action)
		err := s.selectAppendExecItem(ei, invokeTimeout)
		if err != nil {
			freeExecuteItem(ei)
			return nil, err
		}

		select {
		case ret := <-ei.retChan:
			freeExecuteItem(ei)
			return ret.result, ret.err
		case <-syncCallTimer.After(invokeTimeout):
			ei.timeout.Set(true)
			return nil, newSkeletonError(s).SetTimeout()
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				ei.timeout.Set(true)
				return nil, newSkeletonError(s).SetTimeout()
			} else {
				return nil, ctx.Err()
			}
		}
	}
}

func tryDetectRing(current, from unsafe.Pointer, cnt int) (bool, []string) {
	//curr,_ := gid2SkeletonMap.Load(current)
	//fmt.Println(fmt.Sprintf("==== current:%s %p from: %p",curr.(Skeleton).UserData(),current,from))
	if current == gid.Invalid {
		// current not running
		return false, nil
	}
	var callStack []string
	for i := 0; i < cnt; i++ {
		if c, ok := gid2SkeletonMap.Load(from); ok {
			callStack = append(callStack, c.(*skeleton).tag)
			from = c.(*skeleton).getActionFrom()
			//fmt.Println(fmt.Sprintf("=============== from: %p",from))
			if from == current {
				if currentSkeleton, ok := gid2SkeletonMap.Load(current); ok {
					callStack = append(callStack, currentSkeleton.(*skeleton).tag)
				}
				//fmt.Println("is loop call")
				return true, callStack
			}
		} else {
			return false, callStack
		}
	}
	return false, callStack
}

func (s *skeleton) SyncInvoke(fn Action) (interface{}, error) {
	return s.SyncInvokeWithContext(context.Background(), func(ctx context.Context) (interface{}, error) {
		return fn()
	})
}
func (s *skeleton) AsyncInvoke(fn ActionNoReturn) error {
	return s.AsyncInvokeWithContext(context.Background(), func(ctx context.Context) { fn() })
}

func (s *skeleton) SyncInvokeWithContext(ctx context.Context, fn ActionWithContext) (interface{}, error) {
	return s.call(ctx, fn)
}

func (s *skeleton) AsyncInvokeWithContext(ctx context.Context, fn ActionNoReturnWithContext) error {
	ei := allocExecuteItem(ctx, func(ctxUsing context.Context) (interface{}, error) {
		fn(ctxUsing)
		return nil, nil
	})
	err := s.selectAppendExecItem(ei, s.spec.InvokeTimeout)
	if err != nil {
		freeExecuteItem(ei)
	}
	return err
}

var gid2SkeletonMap sync.Map
var gSkeletonID sync2.AtomicUint64
var skeletonID2SkeletonMap sync.Map

// GetRunningSkeleton get running skeleton by skeleton id
func GetRunningSkeleton(skeletonId uint64) (Skeleton, bool) {
	if ss, ok := skeletonID2SkeletonMap.Load(skeletonId); ok && ss.(Skeleton).IsRunning() {
		return ss.(Skeleton), true
	}
	return nil, false
}

var _, _ = GetRunningSkeleton(0)
var _ = GetCurrentSkeleton()

func GetCurrentSkeleton() Skeleton {
	if c, ok := gid2SkeletonMap.Load(gid.Get()); ok {
		return c.(Skeleton)
	}
	return nil
}

type skeletonImpl interface {
	name() string
	exec(*execItem)
	getExecChan() chan *execItem
	getRingExecChan() chan *execItem

	getActionFrom() unsafe.Pointer
	setActionFrom(unsafe.Pointer)
}

type skeleton struct {
	execChan     chan *execItem
	ringExecChan chan *execItem

	execChanMutex sync.RWMutex

	gid        gid.AtomicGid
	skeletonID sync2.AtomicUint64

	*time2.Dispatcher

	spec             *Options
	isRunning        sync2.AtomicInt32
	terminated       sync2.AtomicBool
	stopNotifyChan   chan struct{} // v3通知结束
	stoppedChan      chan struct{} // 供外部监听者
	terminatedChan   chan struct{} // 供外部监听者
	stoppedChanInner chan struct{} // 通知同步等待结束成功
	startedChan      chan struct{} // 通知同步等待启动成功

	// for suture
	sutureToken      suturev4.ServiceToken
	sutureTokenMutex sync.Mutex
	suturev4.Service

	actionFrom unsafe.Pointer

	eventChan chan struct{}
	eventHook func(Skeleton)

	userData interface{}
	_name    string
	tag      string

	lastGLSState map[interface{}]interface{}
	startCount   int
}

func NewSkeleton(tag string, opts ...Option) Skeleton {
	return NewSkeletonWithTag(tag, NewOptions(opts...))
}

func NewSkeletonWithTag(tag string, cc *Options) Skeleton {
	s := &skeleton{spec: cc, tag: tag}
	s.execChan = make(chan *execItem, s.spec.InvokeQueueLen)
	s.ringExecChan = make(chan *execItem, 1)
	s.Dispatcher = time2.NewDispatcher(s.spec.TimerDispatcherLen)
	s.Service = s
	// supervisor server时需要监听,不能再Server中再去创建
	s.startedChan = make(chan struct{}, 1)
	// 逻辑层可能会监听
	s.stoppedChan = make(chan struct{}, 1)
	s.terminatedChan = make(chan struct{}, 1)
	return s
}
func (s *skeleton) RestartCount() int               { return s.startCount }
func (s *skeleton) getExecChan() chan *execItem     { return s.execChan }
func (s *skeleton) getRingExecChan() chan *execItem { return s.ringExecChan }
func (s *skeleton) name() string                    { return s._name }
func (s *skeleton) StoppedChan() chan struct{}      { return s.stoppedChan }
func (s *skeleton) TerminatedChan() chan struct{}   { return s.terminatedChan }
func (s *skeleton) SetUserData(data interface{})    { s.userData = data }
func (s *skeleton) UserData() interface{}           { return s.userData }
func (s *skeleton) getActionFrom() unsafe.Pointer   { return atomic.LoadPointer(&s.actionFrom) }
func (s *skeleton) setActionFrom(p unsafe.Pointer)  { atomic.StorePointer(&s.actionFrom, p) }
func (s *skeleton) GID() gid.Type                   { return s.gid.Get() }
func (s *skeleton) SkeletonID() uint64              { return s.skeletonID.Get() }
func (s *skeleton) IsRunning() bool                 { return s.isRunning.Get() == 1 }
func (s *skeleton) Stop()                           { close(s.stopNotifyChan) } // Stop call from third goroutine, do not call this method
func (s *skeleton) SetEventHook(cc chan struct{}, cb func(Skeleton)) {
	s.eventChan = cc
	s.eventHook = cb
}

func (s *skeleton) inheritGLS() {
	if !s.spec.EnableGls {
		return
	}
	lastState := gls.GetGls(s.gid.Get())
	if lastState == nil {
		lastState = map[interface{}]interface{}{}
	}
	lastState["gid"] = s.gid
	gls.ResetGls(s.gid.Get(), lastState)
}

func (s *skeleton) inheritSkeletonID() {
	// skeleton id在skeleton发生重启时不变化，便于逻辑层缓存
	if s.skeletonID.Get() == 0 {
		s.skeletonID.Set(gSkeletonID.Add(1))
	}
	skeletonID2SkeletonMap.Store(s.skeletonID.Get(), s)
}

// Serve auto called by suture
func (s *skeleton) Serve(ctx context.Context) (err error) {
	s.terminated.Set(false)
	s.gid.Set(gid.Get())
	gid2SkeletonMap.Store(s.gid.Get(), s)

	// 继承gls信息
	s.inheritGLS()

	// 继承skeleton信息
	s.inheritSkeletonID()

	s._name = fmt.Sprintf("skeleton#%p#%d#%s", s.gid.Get(), s.skeletonID.Get(), s.tag)
	// 遗留逻辑，目前不通过stopNotifyChan关闭
	s.stopNotifyChan = make(chan struct{})
	s.stoppedChanInner = make(chan struct{}, 1)

	var notify liveness.Notify
	var liveToken liveness.ServiceToken

	normalStop := false
	defer func() {
		if err == nil {
			if r := recover(); r != nil {
				err = errors.New(r.(string))
			}
		}
		if s.spec.ShouldRestart != nil {
			if should := s.spec.ShouldRestart(s, err.Error()); !should {
				err = suturev4.ErrDoNotRestart
				normalStop = true
			}
		}
		// 重启后gid会变 需要重新设置映射关系
		gid2SkeletonMap.Delete(s.gid.Get())
		if s.spec.LiveProbe {
			liveness.Default.Remove(liveToken)
		}
		if normalStop {
			s.onStop(DefaultTerminateStopTimeout)
		} else {
			s.onTerminate(DefaultTerminateStopTimeout)
		}
	}()
	s.isRunning.Set(1)

	// todo 这部分逻辑考虑重构，独立在skeleton外
	s.startCount++
	log.Info().Str("skeleton", s._name).Int("count", s.startCount).Msg("skeleton start")

	var notifyChan <-chan struct{}
	if s.spec.LiveProbe {
		notifyChan = liveTimeWheel.After(s.spec.LiveNotifyInternal)
		liveToken, notify = liveness.Default.Add(s.name(), int64(s.spec.LiveTTL.Seconds()))
	}

	// notify started
	select {
	case <-s.startedChan:
	default:
	}
	s.startedChan <- struct{}{}

	for {
		select {
		case <-s.eventChan: //block if eventChan is nil
			s.eventChan = nil
			s.eventHook(s)
		case <-s.stopNotifyChan: // v3版本关闭
		case <-ctx.Done(): // v4版本关闭
			normalStop = true
			err = suturev4.ErrDoNotRestart
			return
		case ci := <-s.execChan:
			s.exec(ci)
		case t := <-s.Dispatcher.ChanTimer:
			if t != nil {
				t.Cb()
			}
		case <-notifyChan:
			notify.Notify()
			notifyChan = liveTimeWheel.After(s.spec.LiveNotifyInternal)
		}
	}
}

func (s *skeleton) onTerminate(timeout time.Duration) {
	log.Warn().Str("skeleton", s.name()).Msg("skeleton abnormal terminate")
	if s.spec.OnTerminated != nil {
		s.spec.OnTerminated(s)
	}
	// clean the request queue
	if s.spec.TerminateCleanQueue {
		s.cleanQueue(timeout)
	}
	s.terminated.Set(true)
	select {
	case <-s.terminatedChan:
	default:
	}
	s.terminatedChan <- struct{}{}
}

func (s *skeleton) cleanQueue(timeout time.Duration) {
	// sync to wait service started
	var timeoutC <-chan time.Time
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		timeoutC = timer.C
	}
	shouldContinue := true
	// 依赖于default将chan消耗完毕，存在风险，但是running进入了close状态, chan中不会有更多的请求进入
	for shouldContinue {
		select {
		case ei := <-s.execChan:
			ei.retChan <- &result{result: nil, err: fmt.Errorf("chan req server closed")}
		case ei := <-s.ringExecChan:
			ei.retChan <- &result{result: nil, err: fmt.Errorf("chan req server closed")}
		case <-timeoutC:
			log.Warn().Str("skeleton", s.name()).Dur("timeout", timeout).Int("queue_left", len(s.execChan)).Int("queue_ring_left", len(s.ringExecChan)).Msg("skeleton clean queue with timeout")
		default:
			shouldContinue = false
		}
	}
}

// 正常关闭状态下running设置为0，如果是异常关闭，terminate则仍认为是在运行状态，队列不做清空操作，等待重启后继续运行
func (s *skeleton) onStop(timeout time.Duration) {
	if s.isRunning.CompareAndSwap(1, 0) {
		log.Debug().Str("skeleton", s.name()).Msg("skeleton normal stop")
		s.Dispatcher.Close()

		// clean the request queue
		s.cleanQueue(timeout)

		// 正常结束时删除gls信息，否则gls信息无法实现重启skeleton继承
		if s.spec.EnableGls {
			defer gls.DeleteGls(s.gid.Get())
		}
		skeletonID2SkeletonMap.Delete(s.skeletonID.Get())

		// 通知内部SupervisorStop结束
		select {
		case <-s.stoppedChanInner:
		default:
		}
		s.stoppedChanInner <- struct{}{}

		// 通知外部监听者
		select {
		case <-s.stoppedChan:
		default:
		}
		s.stoppedChan <- struct{}{}
		if s.spec.OnStopped != nil {
			s.spec.OnStopped(s)
		}
	}
}

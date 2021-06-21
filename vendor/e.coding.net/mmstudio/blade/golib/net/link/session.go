package link

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/golib/gpool"
	"e.coding.net/mmstudio/blade/golib/paniccatcher"
	"e.coding.net/mmstudio/blade/golib/sync2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var SessionClosedError = errors.New("session has closed")
var SessionBlockedError = errors.New("session blocked")

var globalSessionId sync2.AtomicUint64

type sessionKeyType struct{}

func (*sessionKeyType) String() string {
	return "session—key-bl9hi867"
}

var SessionKeyForContext = sessionKeyType(struct{}{})

// NewContext creates a new context with session information attached.
func NewContextWithSession(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, SessionKeyForContext, s)
}

// FromContext returns the session information in ctx if it exists.
func SessionFromContext(ctx context.Context) (s Session, ok bool) {
	s, ok = ctx.Value(SessionKeyForContext).(Session)
	return
}

type Session interface {
	ID() uint64
	Transporter() Transporter
	ReceiveLoopWithDelegate(SessionDelegate)
	ReadFrame() (interface{}, error)
	WriteFrame(interface{}) error
	Flush(timeout time.Duration) error
	Close(closeInfo ...string) error
	IsClosed() bool
	ElapsedTime() time.Duration

	// call back and user data
	AddCloseCallback(handler, key interface{}, callback func())
	RemoveCloseCallback(handler, key interface{})

	AddInitCallback(callback func())
	InvokeAndRemoveInitCallback()

	SetUserData(interface{})
	UserData() interface{}

	SetCloseAfterResp(bool)
}

type _session struct {
	id      uint64
	manager *Manager
	//close
	closeChan          chan struct{}
	closeFlag          sync2.AtomicInt32
	closeCallback      callbacks
	initCallback       callbacks
	closeCallbackMutex sync.RWMutex
	initCallbackMutex  sync.RWMutex
	closeAfterResp     sync2.AtomicBool
	// send
	sendChan      chan interface{}
	flushChan     chan *flushContext
	sendChanMutex sync.RWMutex
	// use data bind to each session
	userData interface{}
	// just for read write frame
	transporter Transporter
	// elapsed time
	timeCreate time.Time

	delegate SessionDelegate

	spec *Spec
	gp   *gpool.Pool

	recvMutex sync.Mutex

	// keepalive
	timeLastFrame time.Time
	numInvoke     sync2.AtomicInt64

	baseLogFields  *zerolog.Event
	receivedFrames sync2.AtomicUint64
	sendFrames     sync2.AtomicUint64

	// 以下数据修改通过RWMutex
	closeInfoMutex sync.RWMutex
	closeLogFunc   *zerolog.Event
	closeInfos     []string
	closeFields    *zerolog.Event
}

func NewSession(transporter Transporter, spec *Spec) *_session {
	t := &_session{
		id:          globalSessionId.Add(1),
		closeChan:   make(chan struct{}),
		flushChan:   make(chan *flushContext, 1),
		transporter: transporter,
		spec:        spec,
		timeCreate:  time.Now(),
	}

	t.closeLogFunc = log.Debug()
	t.baseLogFields = log.Info().Uint64("session_id", t.id).Str("remote_address", transporter.RemoteAddr().String()).Str("local_address", transporter.LocalAddr().String())

	t.log(log.Debug(), "session_new")

	if spec.SendChanSize > 0 {
		t.sendChan = make(chan interface{}, spec.SendChanSize)
		go t.sendLoop()
	}

	return t
}

func (s *_session) log(event *zerolog.Event, message string) {
	event.Dur("elapsed", s.ElapsedTime()).Uint64("frames_received", s.receivedFrames.Get()).Uint64("frames_send", s.sendFrames.Get()).Msg(message)
}

func (s *_session) Transporter() Transporter {
	return s.transporter
}

func (s *_session) ID() uint64 {
	return s.id
}

func (s *_session) SetManager(m *Manager) {
	s.manager = m
}

func (s *_session) AddInitCallback(callback func()) {
	s.initCallbackMutex.Lock()
	s.initCallback.Add(s, s, callback)
	s.initCallbackMutex.Unlock()
}
func (s *_session) InvokeAndRemoveInitCallback() {
	s.initCallbackMutex.Lock()
	s.initCallback.Invoke()
	s.initCallback.Remove(s, s)
	s.initCallbackMutex.Unlock()
}

func (s *_session) AddCloseCallback(handler, key interface{}, callback func()) {
	if s.IsClosed() {
		return
	}
	s.closeCallbackMutex.Lock()
	s.closeCallback.Add(handler, key, callback)
	s.closeCallbackMutex.Unlock()
}

func (s *_session) RemoveCloseCallback(handler, key interface{}) {
	if s.IsClosed() {
		return
	}
	s.closeCallbackMutex.Lock()
	s.closeCallback.Remove(handler, key)
	s.closeCallbackMutex.Unlock()
}

func (s *_session) invokeCloseCallbacks() {
	s.closeCallbackMutex.RLock()
	s.closeCallback.Invoke()
	s.closeCallbackMutex.RUnlock()
}

func (s *_session) SetCloseAfterResp(close bool) {
	s.closeAfterResp.Set(close)
}

func (s *_session) SetUserData(userData interface{}) {
	s.userData = userData
}

func (s *_session) UserData() interface{} {
	return s.userData
}

func (s *_session) ElapsedTime() time.Duration {
	if s.timeCreate.IsZero() {
		return 0
	}
	now := time.Now()
	if s.timeCreate.After(now) {
		return s.timeCreate.Sub(now)
	}
	return now.Sub(s.timeCreate)
}

func (s *_session) frameHandler(frame interface{}) {
	s.numInvoke.Add(1)
	handler := func() {
		ctx := NewContextWithSession(context.Background(), s)
		rsp := s.delegate.Invoke(ctx, s, frame)
		if rsp == nil {
			s.log(log.Warn().Bytes("req_bytes", frame.([]byte)), "request got nil response")
			return
		}
		if err := s.WriteFrame(rsp); err != nil {
			s.log(log.Info().Err(err), "WriteFrame with error")
			if owf, ok := s.delegate.(OnWriteFrameError); ok {
				owf.OnWriteFrameError(s, err)
			}
		}
		if s.closeAfterResp.Get() {
			_ = s.Flush(1 * time.Second)
			_ = s.Close("close_after_resp")
		}
		s.numInvoke.Add(-1)
	}
	switch {
	case s.spec.InvokerPerSession < 0:
		// 单用户session模式下使用
		handler()
	case s.spec.InvokerPerSession == 0:
		// rpc session
		go handler()
	default:
		// rpc session
		if s.gp == nil {
			s.gp = gpool.NewPool(s.spec.InvokerPerSession, s.spec.QueuePerInvoker)
		}
		// todo: timeout support under goroutine pool model
		s.gp.JobQueue <- handler
	}
}

func (s *_session) ReceiveLoopWithDelegate(delegate SessionDelegate) {
	s.recvMutex.Lock()
	defer func() {
		s.recvMutex.Unlock()
		_ = s.Close()
	}()
	s.delegate = delegate
	s.timeLastFrame = time.Now()
	orf, _ := s.delegate.(OnReadFrameError)
	for {
		// todo readtimeout效率问题，sysfd是否可以解决？
		if s.spec.ReadTimeout != 0 {
			err := s.transporter.SetReadDeadline(time.Now().Add(s.spec.ReadTimeout))
			if err != nil {
				s.log(log.Debug().Err(err), "link session SetReadDeadline got error")
			}
		}
		// todo 批量读取网络层数据,减少read调用
		if raw, err := s.transporter.Receive(); err != nil {
			if s.spec.IdleTimeout > 0 && s.timeLastFrame.Add(s.spec.IdleTimeout).Before(time.Now()) {
				s.upClosesInfo(log.Info().Err(err), "receive_timeout")
				return
			}
			if isTemporaryErr(err) {
				s.log(log.Debug(), "temporary_err")
				continue
			}
			if err == io.EOF {
				s.upClosesInfo(log.Debug(), "receive_eof")
				return
			}
			if orf != nil {
				orf.OnReadFrameError(s, err)
			} else {
				s.upClosesInfo(log.Error().Err(err), "receive_error")
				return
			}
		} else {
			s.receivedFrames.Add(1)
			s.timeLastFrame = time.Now()
			s.frameHandler(raw)
		}
	}
}

func (s *_session) ReadFrame() (interface{}, error) {
	s.recvMutex.Lock()
	defer s.recvMutex.Unlock()
	ret, err := s.transporter.Receive()
	if err != nil {
		s.upClosesInfo(log.Warn().Err(err), "read_frame_with_err")
		_ = s.Close()
	} else {
		s.receivedFrames.Add(1)
	}
	return ret, err
}

func (s *_session) WriteFrame(b interface{}) error {
	if s.sendChan == nil {
		if s.IsClosed() {
			return SessionClosedError
		}
		s.sendChanMutex.Lock()
		defer s.sendChanMutex.Unlock()
		s.sendFrames.Add(1)
		err := s.transporter.Send(b)
		if err != nil {
			s.upClosesInfo(log.Error().Err(err), "transporter_sync_send_frame_error")
			_ = s.Close()
		}
		return err
	}

	// avoid sendChan closed
	s.sendChanMutex.RLock()
	if s.IsClosed() {
		s.sendChanMutex.RUnlock()
		return SessionClosedError
	}

	// now 'isClosed' is false，we can make sure sendChan is not closed now
	select {
	case s.sendChan <- b:
		s.sendChanMutex.RUnlock()
		return nil
	default:
		s.sendChanMutex.RUnlock()
		s.upClosesInfo(log.Error(), "session_blocked_error")
		_ = s.Close()
		return SessionBlockedError
	}
}

func (s *_session) IsClosed() bool {
	return s.closeFlag.Get() == 1
}

func (s *_session) upClosesInfo(event *zerolog.Event, info string) {
	if s.closeFlag.Get() == 1 {
		return
	}

	s.closeInfoMutex.Lock()
	defer s.closeInfoMutex.Unlock()

	s.closeLogFunc = event
	s.closeInfos = append(s.closeInfos, info)
}

func (s *_session) Close(closeInfo ...string) error {
	if s.closeFlag.CompareAndSwap(0, 1) {

		s.closeInfoMutex.RLock()
		s.log(s.closeLogFunc.Strs("close_info", append(s.closeInfos, closeInfo...)), "session_close")
		s.closeInfoMutex.RUnlock()

		close(s.closeChan)
		if s.sendChan != nil {
			s.sendChanMutex.Lock()
			close(s.sendChan)
			// clean send chan
			if clear, ok := s.transporter.(ClearSendChan); ok {
				clear.ClearSendChan(s, s.sendChan)
			}
			s.sendChanMutex.Unlock()
		}
		_ = s.transporter.Close()

		go func() {
			paniccatcher.Do(func() {
				if s.manager != nil {
					s.manager.delSession(s)
				}
				if s.gp != nil {
					s.gp.Release()
				}
				s.invokeCloseCallbacks()
			}, func(p *paniccatcher.Panic) {
				s.log(log.Error().Interface("reason", p.Reason), "session panic when async callback")
			})
		}()
		return nil
	}
	return SessionClosedError
}

type flushContext struct {
	doneCh chan struct{}
	errCh  chan error
}

// flush will send all data out and close the transport
func (s *_session) Flush(timeout time.Duration) error {
	select {
	case <-s.flushChan:
	default:
	}
	fc := &flushContext{
		doneCh: make(chan struct{}),
		errCh:  make(chan error),
	}
	s.flushChan <- fc

	select {
	case <-s.closeChan:
		return nil
	case <-fc.doneCh:
		return nil
	case err := <-fc.errCh:
		return err
	case <-s.spec.timingWheelForFlushInner.After(timeout):
		return fmt.Errorf("flush error with timeout max %d seconds", timeout/time.Second)
	}
}

func (s *_session) sendLoop() {
	defer func() { _ = s.Close() }()
	for {
		select {
		case fc := <-s.flushChan:
			if err := s.WriteFrame(fc); err != nil {
				fc.errCh <- err
				return
			}
			for b := range s.sendChan {
				if b == nil {
					fc.doneCh <- struct{}{}
					return
				}
				if _, ok := b.(*flushContext); ok {
					fc.doneCh <- struct{}{}
					break
				}
				if err := s.transporter.Send(b); err != nil {
					fc.errCh <- err
					return
				}
				s.sendFrames.Add(1)
			}
		case b, ok := <-s.sendChan:
			if !ok || b == nil {
				s.upClosesInfo(log.Debug().Bool("send_chan_ok", ok), "send_chan_exit")
				return
			}
			if err := s.transporter.Send(b); err != nil {
				s.upClosesInfo(log.Error().Err(err), "transport_async_send_error")
				return
			}
			s.sendFrames.Add(1)
		case <-s.closeChan:
			return
		}
	}
}

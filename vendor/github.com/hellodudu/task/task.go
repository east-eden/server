package task

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"go.uber.org/atomic"
)

var (
	TaskDefaultChannelSize    = 64                       // task channel buffer size
	TaskDefaultExecuteTimeout = time.Second * 5          // execute timeout
	TaskDefaultTimeout        = time.Hour * 24 * 30 * 12 // default timeout
	TaskDefaultSleep          = time.Millisecond * 500   // sleep time 500ms
	ErrTimeout                = errors.New("time out")
)

type TaskHandler func(context.Context, ...interface{}) error
type Task struct {
	c context.Context // control run timeout
	f TaskHandler     // handle function
	e chan<- error    // error returned
	p []interface{}   // params in order
}

type Tasker struct {
	opts     *TaskerOptions
	tasks    chan *Task
	stopChan chan struct{}
	stopOnce *sync.Once
	running  atomic.Bool
}

func NewTasker() *Tasker {
	return &Tasker{
		opts:     &TaskerOptions{},
		tasks:    make(chan *Task, TaskDefaultChannelSize),
		stopChan: make(chan struct{}, 1),
		stopOnce: new(sync.Once),
	}
}

func (t *Tasker) Init(opts ...TaskerOption) {
	t.opts = defaultTaskerOptions()

	for _, o := range opts {
		o(t.opts)
	}
}

func (t *Tasker) ResetTimer() {
	tm := t.opts.timer
	if tm != nil && !tm.Stop() {
		<-tm.C
	}
	tm.Reset(t.opts.d)
}

func (t *Tasker) IsRunning() bool {
	return t.running.Load()
}

func (t *Tasker) AddWait(ctx context.Context, f TaskHandler, p ...interface{}) error {
	subCtx, cancel := context.WithTimeout(ctx, TaskDefaultExecuteTimeout)
	defer cancel()

	e := make(chan error, 1)
	tk := &Task{
		c: subCtx,
		f: f,
		e: e,
		p: make([]interface{}, 0, len(p)),
	}
	tk.p = append(tk.p, p...)
	t.tasks <- tk

	select {
	case err := <-e:
		if err == nil {
			return nil
		}
		return fmt.Errorf("task add with error:%w, chan buff size:%d", err, len(t.tasks))
	case <-subCtx.Done():
		if subCtx.Err() == nil {
			return nil
		}
		return fmt.Errorf("task add with timeout:%w, chan buff size:%d", subCtx.Err(), len(t.tasks))
	}
}

func (t *Tasker) Add(ctx context.Context, f TaskHandler, p ...interface{}) {
	tk := &Task{
		c: ctx,
		f: f,
		e: nil,
		p: make([]interface{}, 0, len(p)),
	}
	tk.p = append(tk.p, p...)
	t.tasks <- tk
}

func (t *Tasker) Run(ctx context.Context) error {
	t.ResetTimer()
	t.running.Store(true)

	defer func() {
		if err := recover(); err != nil {
			stack := string(debug.Stack())
			fmt.Printf("catch exception:%v, panic recovered with stack:%s", err, stack)
		}
		t.stop()
	}()

	if len(t.opts.startFns) > 0 {
		for _, fn := range t.opts.startFns {
			fn()
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case h, ok := <-t.tasks:
			if !ok {
				return nil
			} else {
				select {
				case <-h.c.Done():
					continue
				default:
				}

				err := h.f(h.c, h.p...)
				if h.e != nil {
					h.e <- err // handle result
				}

				if err == nil {
					continue
				}
			}

		case <-t.opts.timer.C:
			return ErrTimeout

		default:
			now := time.Now()
			if t.opts.updateFn != nil {
				t.opts.updateFn() // update callback
			}
			d := time.Since(now)
			time.Sleep(t.opts.sleep - d)
		}
	}
}

func (t *Tasker) Stop() {
	t.stopOnce.Do(func() {
		close(t.tasks)
		<-t.stopChan
	})
}

func (t *Tasker) stop() {
	if len(t.opts.stopFns) > 0 {
		for _, fn := range t.opts.stopFns {
			fn()
		}
	}

	t.opts.timer.Stop()
	t.running.Store(false)
	close(t.stopChan)
}

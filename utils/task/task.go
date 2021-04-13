package task

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/funplus/server/utils"
)

var (
	TaskDefaultExecuteTimeout = time.Second * 5          // execute 执行超时
	TaskDefaultTimeout        = time.Hour * 24 * 30 * 12 // 默认超时时间
	TaskDefaultSleep          = time.Millisecond * 100   // 默认睡眠100ms
	ErrTimeout                = errors.New("time out")
)

type ContextDoneHandler func()
type ErrorFilter func(error) error
type DefaultUpdate func()

type TaskHandler func(context.Context, ...interface{}) error
type Task struct {
	C context.Context // 超时控制
	F TaskHandler     // 任务处理函数
	E chan<- error    // 返回error
	P []interface{}   // 任务参数
}

type TaskerOption func(*TaskerOptions)
type TaskerOptions struct {
	cdh ContextDoneHandler // context done回调
	tm  *time.Timer        // 超时处理
	d   time.Duration      // 超时时间
	du  DefaultUpdate      // default处理
	sd  time.Duration      // 默认睡眠时间
}

type Tasker struct {
	opts  *TaskerOptions
	tasks chan *Task
}

func NewTasker(max int32) *Tasker {
	return &Tasker{
		opts:  &TaskerOptions{},
		tasks: make(chan *Task, max),
	}
}

func DefaultTaskerOptions() *TaskerOptions {
	return &TaskerOptions{
		d:   TaskDefaultTimeout,
		cdh: nil,
		tm:  time.NewTimer(TaskDefaultTimeout),
		du:  nil,
		sd:  TaskDefaultSleep,
	}
}

func WithContextDoneHandler(h ContextDoneHandler) TaskerOption {
	return func(o *TaskerOptions) {
		o.cdh = h
	}
}

func WithTimeout(d time.Duration) TaskerOption {
	return func(o *TaskerOptions) {
		o.d = d
	}
}

func WithDefaultUpdate(du DefaultUpdate) TaskerOption {
	return func(o *TaskerOptions) {
		o.du = du
	}
}

func (t *Tasker) Init(opts ...TaskerOption) {
	t.opts = DefaultTaskerOptions()

	for _, o := range opts {
		o(t.opts)
	}
}

func (t *Tasker) ResetTimer() {
	tm := t.opts.tm
	if tm != nil && !tm.Stop() {
		<-tm.C
	}
	tm.Reset(t.opts.d)
	// t.opts.tm.Reset(t.opts.d)
}

func (t *Tasker) Add(c context.Context, f TaskHandler, p ...interface{}) error {
	subCtx, cancel := utils.WithTimeoutContext(c, TaskDefaultExecuteTimeout)
	defer cancel()

	e := make(chan error, 1)
	tk := &Task{
		C: subCtx,
		F: f,
		E: e,
		P: make([]interface{}, 0, len(p)),
	}
	tk.P = append(tk.P, p...)
	t.tasks <- tk

	select {
	case err := <-e:
		return err
	case <-subCtx.Done():
		return subCtx.Err()
	}
}

func (t *Tasker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			if t.opts.cdh != nil {
				t.opts.cdh()
			}
			return nil

		case h, ok := <-t.tasks:
			if !ok {
				return nil
			} else {
				err := h.F(h.C, h.P...)
				h.E <- err
				if err == nil {
					continue
				}
			}

		case <-t.opts.tm.C:
			return ErrTimeout

		default:
			now := time.Now()
			if t.opts.du != nil {
				t.opts.du()
			}
			d := time.Since(now)
			time.Sleep(t.opts.sd - d)
		}
	}
}

func (t *Tasker) Stop() {
	close(t.tasks)
	t.opts.tm.Stop()
}

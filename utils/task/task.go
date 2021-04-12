package task

import (
	"context"
	"time"

	"bitbucket.org/funplus/server/utils"
)

var (
	TaskDefaultTimeout = time.Second * 5
)

type TaskHandler func(context.Context, ...interface{}) error
type Task struct {
	C context.Context // 超时控制
	F TaskHandler     // 任务处理函数
	E chan<- error    // 返回error
	P []interface{}   // 任务参数
}

type Tasker struct {
	tasks chan *Task
}

func NewTasker(max int32) *Tasker {
	return &Tasker{
		tasks: make(chan *Task, max),
	}
}

func (t *Tasker) Execute(c context.Context, f TaskHandler, p ...interface{}) error {
	subCtx, cancel := utils.WithTimeoutContext(c, TaskDefaultTimeout)
	defer cancel()

	e := make(chan error, 1)
	tk := &Task{
		C: c,
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
			return nil

		case h, ok := <-t.tasks:
			if !ok {
				return nil
			} else {
				err := h.F(h.C, h.P...)
				h.E <- err
				if err != nil {
					return err
				}
			}
		}
	}
}

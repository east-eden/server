package task

import "time"

type StartFn func()
type StopFn func()
type UpdateFn func()

type TaskerOption func(*TaskerOptions)
type TaskerOptions struct {
	startFns []StartFn // start callback
	stopFn   StopFn    // task stop callback
	updateFn UpdateFn  // default update callback
	timer    *time.Timer
	d        time.Duration // timeout duration
	sleep    time.Duration // sleep duration
}

func defaultTaskerOptions() *TaskerOptions {
	return &TaskerOptions{
		d:        TaskDefaultTimeout,
		startFns: make([]StartFn, 0, 5),
		stopFn:   nil,
		updateFn: nil,
		timer:    time.NewTimer(TaskDefaultTimeout),
		sleep:    TaskDefaultSleep,
	}
}

func WithStartFns(f ...StartFn) TaskerOption {
	return func(o *TaskerOptions) {
		o.startFns = o.startFns[:0]
		o.startFns = append(o.startFns, f...)
	}
}

func WithStopFn(f StopFn) TaskerOption {
	return func(o *TaskerOptions) {
		o.stopFn = f
	}
}

func WithUpdateFn(f UpdateFn) TaskerOption {
	return func(o *TaskerOptions) {
		o.updateFn = f
	}
}

func WithTimeout(d time.Duration) TaskerOption {
	return func(o *TaskerOptions) {
		o.d = d
	}
}

func WithSleep(d time.Duration) TaskerOption {
	return func(o *TaskerOptions) {
		o.sleep = d
	}
}

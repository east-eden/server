package channelwriter

import (
	"runtime/debug"
	"sync"
	"time"
)

var (
	WriteBufferSize = 1024 // write buffer size
	SleepDuration   = 50 * time.Millisecond
	FlushInterval   = 2 * time.Second
)

type ChannelWriter struct {
	opts               *Options
	datas              []interface{}
	writeChan          chan interface{}
	closeChan          chan bool
	stopChan           chan bool
	flushImmediateChan chan bool
	d                  time.Duration
	t                  *time.Timer
	once               sync.Once
}

func NewChannelWriter(opts ...Option) *ChannelWriter {
	w := &ChannelWriter{
		opts:               defaultOptions(),
		closeChan:          make(chan bool, 1),
		stopChan:           make(chan bool, 1),
		flushImmediateChan: make(chan bool, 1),
	}

	for _, o := range opts {
		o(w.opts)
	}

	w.datas = make([]interface{}, 0, w.opts.writeBufferSize)
	w.writeChan = make(chan interface{}, w.opts.chanBufferSize)
	w.t = time.NewTimer(w.opts.flushInterval)

	w.run()
	return w
}

func (w *ChannelWriter) ResetFlushInterval(d time.Duration) {
	if w.t != nil && !w.t.Stop() {
		<-w.t.C
	}
	w.opts.flushInterval = d
	w.t.Reset(d)
}

func (w *ChannelWriter) Write(data interface{}) {
	w.writeChan <- data
}

func (w *ChannelWriter) Flush() {
	w.flushImmediateChan <- true
}

func (w *ChannelWriter) Stop() {
	w.once.Do(func() {
		close(w.closeChan)
		<-w.stopChan
	})
}

func (w *ChannelWriter) run() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				w.opts.logger.Printf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			w.stopChan <- true
		}()

		for {
			select {
			case <-w.closeChan:
				w.flush()
				return
			case model := <-w.writeChan:
				w.datas = append(w.datas, model)
			case <-w.t.C:
				w.flush()
				w.t.Reset(w.d)
			case <-w.flushImmediateChan:
				w.flush()
			default:
				if len(w.datas) < WriteBufferSize {
					time.Sleep(SleepDuration)
				} else {
					w.flush()
				}
			}
		}
	}()
}

func (w *ChannelWriter) flush() {
	if len(w.datas) <= 0 {
		return
	}

	err := w.opts.flushHandler(w.datas)
	if err != nil {
		w.opts.logger.Printf("flush failed due to %v", err)
		return
	}
	w.datas = w.datas[:0]
}

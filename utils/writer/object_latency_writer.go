package writer

import (
	"sync"
	"time"
)

// maxLatencyWriter from httputil
type ObjectLatencyWriter struct {
	dst     ObjectWriteFlusher
	latency time.Duration // non-zero; negative means to flush immediately

	mu           sync.Mutex // protects t, flushPending, and dst.Flush
	t            *time.Timer
	flushPending bool
}

// io.writer interface
func (w *ObjectLatencyWriter) Write(p interface{}) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	err = w.dst.Write(p)
	if w.latency < 0 {
		w.dst.Flush()
		return
	}
	if w.flushPending {
		return
	}
	if w.t == nil {
		w.t = time.AfterFunc(w.latency, w.delayedFlush)
	} else {
		w.t.Reset(w.latency)
	}
	w.flushPending = true
	return
}

func (w *ObjectLatencyWriter) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.flushPending = false
	if w.t != nil {
		w.t.Stop()
	}
	w.dst.Flush()
}

func (w *ObjectLatencyWriter) delayedFlush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.flushPending { // if stop was called but AfterFunc already started this goroutine
		return
	}
	w.dst.Flush()
	w.flushPending = false
}

func NewObjectLatencyWriter(wf ObjectWriteFlusher, latency time.Duration) ObjectWriter {
	return &ObjectLatencyWriter{
		dst:     wf,
		latency: latency,
	}
}

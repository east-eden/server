package writer

import (
	"io"
	"time"
)

var (
	DefaultWriterLatency = 500 * time.Millisecond
	DefaultWriterSize    = 4096
)

type WriteFlusher interface {
	io.Writer
	Flush() error
}

type Writer interface {
	Write([]byte) (int, error)
	Stop()
}

func NewWriter(wf WriteFlusher, latency time.Duration) Writer {
	return NewLatencyWriter(wf, latency)
}

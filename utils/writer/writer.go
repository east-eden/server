package writer

import (
	"io"
	"time"
)

var (
	DefaultWriterLatency    = 500 * time.Millisecond
	DefaultBinaryWriterSize = 4096
	DefaultObjectWriterSize = 100
)

// binary writer interface
type BinaryWriteFlusher interface {
	io.Writer
	Flush() error
}

type BinaryWriter interface {
	Write([]byte) (int, error)
	Stop()
}

func NewBinaryWriter(wf BinaryWriteFlusher, latency time.Duration) BinaryWriter {
	return NewBinaryLatencyWriter(wf, latency)
}

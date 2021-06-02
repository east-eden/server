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

// object writer interface
type ObjectWriteFlusher interface {
	Write(interface{}) error
	Flush() error
}

type ObjectWriter interface {
	Write(interface{}) error
	Flush() error
	Stop()
}

func NewBinaryWriter(wf BinaryWriteFlusher, latency time.Duration) BinaryWriter {
	return NewBinaryLatencyWriter(wf, latency)
}

func NewObjectWriter(wf ObjectWriteFlusher, latency time.Duration) ObjectWriter {
	return NewObjectLatencyWriter(wf, latency)
}

package writer

import (
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// test cases
var (
	cases = map[string]struct {
		latency  time.Duration
		wait     bool
		inputBuf [][]byte
		flushBuf []byte
	}{
		"small_piece1": {
			latency: time.Millisecond * 200,
			wait:    false,
			inputBuf: [][]byte{
				[]byte("123"),
			},
			flushBuf: []byte(""),
		},

		"small_piece2": {
			latency: time.Millisecond * 200,
			wait:    true,
			inputBuf: [][]byte{
				[]byte("123"),
			},
			flushBuf: []byte("123"),
		},

		"small_piece3": {
			latency: -1,
			wait:    false,
			inputBuf: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("7890"),
				[]byte("abc"),
			},
			flushBuf: []byte("1234567890abc"),
		},

		"big_piece": {
			latency: -1,
			wait:    false,
			inputBuf: [][]byte{
				[]byte("1234567890abcd"),
			},
			flushBuf: []byte("1234567890abcd"),
		},
	}
)

type fakeFlusher struct {
	flush []byte
	buf   []byte
	size  int
	io.Writer
}

func (f *fakeFlusher) Flush() error {
	f.flush = append(f.flush, f.buf...)
	f.buf = f.buf[:0]
	return nil
}

func (f *fakeFlusher) Write(p []byte) (n int, err error) {
	f.buf = append(f.buf, p...)
	if len(f.buf) >= f.size {
		f.Flush()
	}
	return len(p), nil
}

func TestTransportWriter(t *testing.T) {

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			f := &fakeFlusher{
				size:  10,
				buf:   make([]byte, 0),
				flush: make([]byte, 0),
			}
			w := NewBinaryWriter(f, cs.latency)

			for _, input := range cs.inputBuf {
				_, _ = w.Write(input)
			}

			if cs.wait {
				time.Sleep(cs.latency * 2)
			}

			diff := cmp.Diff(f.flush, cs.flushBuf)
			if diff != "" {
				t.Fatalf(diff)
			}

			w.Stop()
		})
	}
}

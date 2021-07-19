package channelwriter

import (
	"io"
	"log"
	"time"
)

type FlushHandler func([]interface{}) error

type Option func(*Options)
type Options struct {
	writeBufferSize int
	chanBufferSize  int
	flushInterval   time.Duration
	sleepDuration   time.Duration
	logger          *log.Logger
	flushHandler    FlushHandler
}

func defaultOptions() *Options {
	return &Options{
		writeBufferSize: 1024,
		chanBufferSize:  64,
		flushInterval:   2 * time.Second,
		sleepDuration:   50 * time.Millisecond,
		logger:          log.Default(),
		flushHandler:    func([]interface{}) error { return nil },
	}
}

func WithWriteBufferSize(sz int) Option {
	return func(o *Options) {
		o.writeBufferSize = sz
	}
}

func WithChanBufferSize(sz int) Option {
	return func(o *Options) {
		o.chanBufferSize = sz
	}
}

func WithFlushInterval(t time.Duration) Option {
	return func(o *Options) {
		o.flushInterval = t
	}
}

func WithSleepDuration(d time.Duration) Option {
	return func(o *Options) {
		o.sleepDuration = d
	}
}

func WithLogger(w io.Writer) Option {
	return func(o *Options) {
		o.logger = log.New(w, "channelwriter: ", log.Lmsgprefix|log.LstdFlags)
		log.SetOutput(w)
	}
}

func WithFlushHandler(h FlushHandler) Option {
	return func(o *Options) {
		o.flushHandler = h
	}
}

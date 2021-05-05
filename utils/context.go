package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// create sub context with father context's deadline * 3/4
func WithTimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	d, ok := ctx.Deadline()
	if !ok {
		return context.WithTimeout(ctx, timeout)
	} else {
		return context.WithTimeout(ctx, time.Until(d)*3/4)
	}
}

func WithSignaledCancel(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	sub, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
	}()

	return sub
}

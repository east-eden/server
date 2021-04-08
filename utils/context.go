package utils

import (
	"context"
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

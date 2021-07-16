package link

import (
	"context"
)

type SessionDelegate interface {
	Invoke(context.Context, Session, interface{}) []byte
}

type SessionDelegateFunc func(context.Context, Session, interface{}) []byte

func (pf SessionDelegateFunc) Invoke(ctx context.Context, sess Session, frame interface{}) []byte {
	return pf(ctx, sess, frame)
}

type OnReadFrameError interface {
	OnReadFrameError(Session, error)
}

type OnWriteFrameError interface {
	OnWriteFrameError(Session, error)
}

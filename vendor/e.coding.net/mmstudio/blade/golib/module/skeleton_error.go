package module

import (
	"fmt"
	"os"
)

var _ = IsSkeletonClosed(nil)
var _ = IsSkeletonTimeout(newSkeletonError(&mockSkeletonInfoForError{}))
var _ = os.IsTimeout(newSkeletonError(&mockSkeletonInfoForError{}))

type SkeletonInfoForError interface {
	IsRunning() bool
	name() string
}

type mockSkeletonInfoForError struct{}

func (mockSkeletonInfoForError) IsRunning() bool { return false }
func (mockSkeletonInfoForError) name() string    { return "" }

type SkeletonPanic interface{ SkeletonPanic() bool }
type SkeletonBlocked interface{ SkeletonBlocked() bool }
type SkeletonClosed interface{ SkeletonClosed() bool }

func IsSkeletonPanic(err error) bool {
	if err == nil {
		return false
	}
	terr, ok := err.(SkeletonPanic)
	return ok && (terr.SkeletonPanic())
}

func IsSkeletonBlocked(err error) bool {
	if err == nil {
		return false
	}
	terr, ok := err.(SkeletonBlocked)
	return ok && (terr.SkeletonBlocked())
}

func IsSkeletonClosed(err error) bool {
	if err == nil {
		return false
	}
	terr, ok := err.(SkeletonClosed)
	return ok && (terr.SkeletonClosed())
}

func IsSkeletonTimeout(err error) bool {
	if err == nil {
		return false
	}
	return os.IsTimeout(err)
}

type SkeletonError struct {
	skeleton string
	reason   []string
	timeout  bool
	blocked  bool
	closed   bool
	panic    bool
}

func (e *SkeletonError) Error() string {
	return fmt.Sprintf("%s error with: %s", e.skeleton, e.reason)
}
func (e *SkeletonError) SkeletonBlocked() bool {
	return e.blocked
}
func (e *SkeletonError) SkeletonClosed() bool {
	return e.closed
}
func (e *SkeletonError) SkeletonPanic() bool {
	return e.panic
}

func (e *SkeletonError) Timeout() bool {
	return e.timeout
}

func (e *SkeletonError) SetTimeout() *SkeletonError {
	e.timeout = true
	e.reason = append(e.reason, "timeout")
	return e
}
func (e *SkeletonError) SetPanic(reason string) *SkeletonError {
	e.panic = true
	e.reason = append(e.reason, "panic", reason)
	return e
}
func (e *SkeletonError) SetBlocked() *SkeletonError {
	e.blocked = true
	e.reason = append(e.reason, "blocked")
	return e
}

func (e *SkeletonError) SetClosed() *SkeletonError {
	e.closed = true
	e.reason = append(e.reason, "closed")
	return e
}

func newSkeletonError(s SkeletonInfoForError) *SkeletonError {
	if s == nil {
		return nil
	}
	e := &SkeletonError{skeleton: s.name()}
	if !s.IsRunning() {
		e.reason = append(e.reason, "closed")
		e.closed = true
	}
	return e
}

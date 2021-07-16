package gid

import (
	"sync/atomic"
	"unsafe"
)

type AtomicGid struct {
	gid atomic.Value
}

func (s *AtomicGid) Set(p unsafe.Pointer) {
	s.gid.Store(p)
}

func (s *AtomicGid) Get() unsafe.Pointer {
	if t := s.gid.Load(); t == nil {
		return nil
	} else {
		return t.(unsafe.Pointer)
	}
}

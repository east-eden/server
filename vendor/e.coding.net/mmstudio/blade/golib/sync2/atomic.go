// Copyright 2013, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync2

import (
	"sync/atomic"
	"time"
)

type AtomicBool struct{ v int32 }

func (b *AtomicBool) Set(v bool) {
	if v {
		atomic.StoreInt32(&b.v, 1)
	} else {
		atomic.StoreInt32(&b.v, 0)
	}
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&b.v) == 1
}

type AtomicInt32 int32

func (i *AtomicInt32) Add(n int32) int32 {
	return atomic.AddInt32((*int32)(i), n)
}

func (i *AtomicInt32) Set(n int32) {
	atomic.StoreInt32((*int32)(i), n)
}

func (i *AtomicInt32) Get() int32 {
	return atomic.LoadInt32((*int32)(i))
}

func (i *AtomicInt32) CompareAndSwap(oldval, newval int32) (swapped bool) {
	return atomic.CompareAndSwapInt32((*int32)(i), oldval, newval)
}

type AtomicUint32 uint32

func (i *AtomicUint32) Add(n uint32) uint32 {
	return atomic.AddUint32((*uint32)(i), n)
}

func (i *AtomicUint32) Set(n uint32) {
	atomic.StoreUint32((*uint32)(i), n)
}

func (i *AtomicUint32) Get() uint32 {
	return atomic.LoadUint32((*uint32)(i))
}

func (i *AtomicUint32) CompareAndSwap(oldval, newval uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32((*uint32)(i), oldval, newval)
}

type AtomicInt64 int64

func (i *AtomicInt64) Add(n int64) int64 {
	return atomic.AddInt64((*int64)(i), n)
}

func (i *AtomicInt64) Set(n int64) {
	atomic.StoreInt64((*int64)(i), n)
}

func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

func (i *AtomicInt64) CompareAndSwap(oldval, newval int64) (swapped bool) {
	return atomic.CompareAndSwapInt64((*int64)(i), oldval, newval)
}

type AtomicUint64 uint64

func (i *AtomicUint64) Add(n uint64) uint64 {
	return atomic.AddUint64((*uint64)(i), n)
}

func (i *AtomicUint64) Set(n uint64) {
	atomic.StoreUint64((*uint64)(i), n)
}

func (i *AtomicUint64) Get() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}

func (i *AtomicUint64) CompareAndSwap(oldval, newval uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64((*uint64)(i), oldval, newval)
}

type AtomicDuration int64

func (d *AtomicDuration) Add(duration time.Duration) time.Duration {
	return time.Duration(atomic.AddInt64((*int64)(d), int64(duration)))
}

func (d *AtomicDuration) Set(duration time.Duration) {
	atomic.StoreInt64((*int64)(d), int64(duration))
}

func (d *AtomicDuration) Get() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

func (d *AtomicDuration) CompareAndSwap(oldval, newval time.Duration) (swapped bool) {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(oldval), int64(newval))
}

// String is an atomic type-safe wrapper for string values.
type AtomicString struct {
	v atomic.Value
}

var _zeroString string

// NewString creates a new String.
func NewAtomicString(v string) *AtomicString {
	x := &AtomicString{}
	if v != _zeroString {
		x.Set(v)
	}
	return x
}

// Load atomically loads the wrapped string.
func (x *AtomicString) Get() string {
	if v := x.v.Load(); v != nil {
		return v.(string)
	}
	return _zeroString
}

// Store atomically stores the passed string.
func (x *AtomicString) Set(v string) {
	x.v.Store(v)
}

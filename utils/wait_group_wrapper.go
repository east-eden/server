package utils

import (
	"fmt"
	"runtime"
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("WaitGroupWrapper: panic recovered: %s\ncall stack: %s\n", r, buf)
			}

			w.Done()
		}()

		cb()
	}()
}

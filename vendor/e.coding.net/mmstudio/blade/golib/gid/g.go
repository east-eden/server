// Package g exposes goroutine struct g to user space.
package gid

import (
	"unsafe"
)

type Type = unsafe.Pointer

var Invalid unsafe.Pointer = nil

func getg() Type

// G returns current g (the goroutine struct) to user space.
func Get() Type {
	return getg()
}

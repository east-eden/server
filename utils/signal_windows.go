//go:build windows

package utils

import "os"

// ShutdownSignals 返回触发优雅退出的信号。
func ShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}

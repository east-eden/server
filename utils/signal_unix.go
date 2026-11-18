//go:build !windows

package utils

import (
	"os"
	"syscall"
)

// ShutdownSignals 返回触发优雅退出的信号。
// SIGINT 用于本地/调试手动中断，SIGTERM 用于 k8s 删除 pod、缩容和滚动更新。
func ShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM}
}

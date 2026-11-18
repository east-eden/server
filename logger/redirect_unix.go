//go:build !windows

package logger

import (
	"fmt"
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
)

var globalErrorFile *os.File

// RedirectStderr 把进程 stderr 重定向到独立文件，
// 避免运行时 panic 的堆栈丢失在容器日志之外。
func RedirectStderr(appName string) {
	errorFn := fmt.Sprintf("data/log/%s_error.log", appName)

	var err error
	globalErrorFile, err = os.OpenFile(errorFn, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Error().Err(err).Str("file", errorFn).Msg("open stderr redirect file failed")
		return
	}

	if err = syscall.Dup2(int(globalErrorFile.Fd()), int(os.Stderr.Fd())); err != nil {
		log.Error().Err(err).Msg("redirect stderr failed")
	}
}

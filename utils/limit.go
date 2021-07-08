package utils

import "syscall"

func Setrlimit(limit uint64) {
	var lim syscall.Rlimit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim)
	if lim.Cur < limit || lim.Max < limit {
		lim.Cur = limit
		lim.Max = limit
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
	}
}

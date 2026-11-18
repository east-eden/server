//go:build windows

package logger

// RedirectStderr 在 windows 上暂不支持重定向 stderr。
func RedirectStderr(appName string) {}

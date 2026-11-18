package utils

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// logWithParams 把可变参数按 p0/p1/... 的形式挂到日志事件上
func logWithParams(event *zerolog.Event, param []any, skipNil bool) *zerolog.Event {
	for k, v := range param {
		if skipNil && v == nil {
			continue
		}

		event = event.Interface(fmt.Sprintf("p%d", k), v)
	}

	return event
}

// ErrCheck print err and return pass false if not nil
func ErrCheck(err error, msg string, param ...any) (pass bool) {
	if err != nil {
		logWithParams(log.Error().CallerSkipFrame(1).Err(err), param, true).Msg(msg)
		return false
	}

	return true
}

// WarnCheck print err and return pass false if not nil
func WarnCheck(err error, msg string, param ...any) (pass bool) {
	if err != nil {
		logWithParams(log.Warn().CallerSkipFrame(1).Err(err), param, true).Msg(msg)
		return false
	}

	return true
}

// InfoCheck print err and return pass false if not nil
func InfoCheck(err error, msg string, param ...any) (pass bool) {
	if err != nil {
		logWithParams(log.Info().CallerSkipFrame(1).Err(err), param, true).Msg(msg)
		return false
	}

	return true
}

// ErrPrint print err if not nil
func ErrPrint(err error, msg string, param ...any) {
	if err != nil {
		logWithParams(log.Warn().CallerSkipFrame(1).Err(err), param, false).Msg(msg)
	}
}

// WarnPrint print err if not nil
func WarnPrint(err error, msg string, param ...any) {
	if err != nil {
		logWithParams(log.Warn().CallerSkipFrame(1).Err(err), param, false).Msg(msg)
	}
}

// InfoPrint print err if not nil
func InfoPrint(err error, msg string, param ...any) {
	if err != nil {
		logWithParams(log.Info().CallerSkipFrame(1).Err(err), param, false).Msg(msg)
	}
}

package utils

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// print err and return pass false if not nil
func ErrCheck(err error, msg string, param ...any) (pass bool) {
	if err != nil {
		event := log.Error().Err(err)
		for k, v := range param {
			event = event.Interface(fmt.Sprintf("p%d", k), v)
		}

		event.Caller(1).Msg(msg)
		pass = false
		return
	}

	pass = true
	return
}

// print err if not nil
func ErrPrint(err error, msg string, param ...any) {
	if err != nil {
		event := log.Warn().Err(err)
		for k, v := range param {
			event = event.Interface(fmt.Sprintf("p%d", k), v)
		}

		event.Caller(1).Msg(msg)
	}
}

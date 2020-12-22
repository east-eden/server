package utils

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func ErrCheck(err error, msg string, param ...interface{}) bool {
	if err != nil {
		event := log.Error().Err(err)
		for k, v := range param {
			event = event.Interface(fmt.Sprintf("p%d", k), v)
		}

		event.Msg(msg)
		return true
	}

	return false
}

func ErrPrint(err error, msg string, param ...interface{}) {
	if err != nil {
		event := log.Error().Err(err)
		for k, v := range param {
			event = event.Interface(fmt.Sprintf("p%d", k), v)
		}

		event.Msg(msg)
	}
}

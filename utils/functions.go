package utils

import (
	"errors"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

var ExceptionErr = errors.New("recover with exception")

func CaptureException() error {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		log.Error().Msgf("panic: Recovered in err", err, stack)
		return ExceptionErr
	}

	return nil
}

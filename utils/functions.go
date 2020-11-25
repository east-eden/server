package utils

import (
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

func CaptureException() {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		log.Error().Msgf("panic: Recovered in err", err, stack)
	}
}

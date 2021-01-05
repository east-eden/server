package utils

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// return event and pass
func ErrCheck(err error, param ...interface{}) (*zerolog.Event, bool) {
	if err != nil {
		event := log.Error().Err(err)
		for k, v := range param {
			event = event.Interface(fmt.Sprintf("p%d", k), v)
		}

		return event, false
	}

	return nil, true
}

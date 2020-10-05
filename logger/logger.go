package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// log file name
	t := time.Now()
	fileTime := fmt.Sprintf("%d-%d-%d %d-%d-%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	logFn := fmt.Sprintf("data/log/gate_%s.log", fileTime)

	file, err := os.OpenFile(logFn, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal().Err(err)
	}

	// set writer
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file))
}

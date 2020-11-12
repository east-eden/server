package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rotate "gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zerolog.Logger

func InitLogger(appName string) {
	// log file name
	t := time.Now()
	fileTime := fmt.Sprintf("%d-%d-%d %d-%d-%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	logFn := fmt.Sprintf("data/log/%s_%s.log", appName, fileTime)

	//file, err := os.OpenFile(logFn, os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//log.Fatal().Err(err)
	//}

	// set console writer and file rotate writer
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, &rotate.Logger{
		Filename:   logFn,
		MaxSize:    200, // megabytes
		MaxBackups: 3,
		MaxAge:     15, //days
	})).With().Caller().Logger()

	Logger = &log.Logger
}

package logger

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rotate "gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zerolog.Logger
var callerPrefixStrim string = "east-eden/server/" // 日志中去除包含此地址的前缀字符

func InitLogger(appName string) {
	// log file name
	t := time.Now()
	fileTime := fmt.Sprintf("%d-%d-%d %d-%d-%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	logFn := fmt.Sprintf("data/log/%s_%s.log", appName, fileTime)

	//file, err := os.OpenFile(logFn, os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//log.Fatal().Err(err)
	//}

	// caller marshal func
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		idx := strings.LastIndex(file, callerPrefixStrim)
		if idx == -1 {
			return file + ":" + strconv.Itoa(line)
		} else {
			return file[idx+len(callerPrefixStrim):] + ":" + strconv.Itoa(line)
		}
	}

	// set console writer and file rotate writer
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, &rotate.Logger{
		Filename:   logFn,
		MaxSize:    200, // megabytes
		MaxBackups: 3,
		MaxAge:     15, //days
	})).With().Logger()

	Logger = &log.Logger
}

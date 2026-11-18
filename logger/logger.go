package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	zero_micro "github.com/asim/go-micro/plugins/logger/zerolog/v4"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	micro_logger "go-micro.dev/v4/logger"
	rotate "gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger        *zerolog.Logger
	ConsoleWriter io.Writer
	RotateWriter  io.Writer
	MultiWriter   io.Writer

	// callerPrefix 日志中需要被裁剪掉的源码路径前缀
	callerPrefix = projectPathName + "/"
)

const (
	projectPathName = "github.com/east-eden/server"

	logMaxSize    = 256 // 单个日志文件上限，单位 MB
	logMaxBackups = 10
	logMaxAge     = 7 // 日志保留天数
)

func init() {
	InitZerolog()
}

// InitZerolog 统一 zerolog 的时间格式与时区。
// 必须在使用任何 log 调用之前执行，因此放在 init 中。
func InitZerolog() {
	zerolog.TimeFieldFormat = strings.Join([]string{time.DateTime, "000"}, ".")
	zerolog.TimestampFunc = utils.NowInLocation
}

// InitLogger 初始化进程日志：控制台 + 文件轮转，并同步 go-micro 的默认 logger。
func InitLogger(appName string) {
	logFn := fmt.Sprintf("data/log/%s.log", appName)

	// caller marshal func：裁剪掉仓库绝对路径，只保留相对路径
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		idx := strings.LastIndex(file, callerPrefix)
		if idx == -1 {
			return file + ":" + cast.ToString(line)
		}

		return file[idx+len(callerPrefix):] + ":" + cast.ToString(line)
	}

	// console writer
	if runtime.GOOS == "windows" {
		ConsoleWriter = &zerolog.ConsoleWriter{TimeFormat: zerolog.TimeFieldFormat, NoColor: true, Out: os.Stdout}
	} else {
		ConsoleWriter = &zerolog.ConsoleWriter{TimeFormat: zerolog.TimeFieldFormat, Out: os.Stdout}
	}

	// rotate writer
	RotateWriter = &rotate.Logger{
		Filename:   logFn,
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAge,
	}

	MultiWriter = io.MultiWriter(ConsoleWriter, RotateWriter)
	log.Logger = log.Output(MultiWriter).With().
		Caller().
		Str("logstore", appName).
		Logger()

	Logger = &log.Logger

	// 让 go-micro 内部日志复用同一套输出
	micro_logger.DefaultLogger = zero_micro.NewLogger(
		micro_logger.WithOutput(MultiWriter),
		micro_logger.WithLevel(micro_logger.Level(Logger.GetLevel()-1)),
		micro_logger.WithCallerSkipCount(zerolog.CallerSkipFrameCount),
	)

	InitZerolog()
}

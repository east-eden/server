package utils

import (
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"strings"

	"github.com/east-eden/server/define"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////
// exception
func CaptureException(p ...interface{}) {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		log.Error().Caller(1).Interface("exception_param", p).Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
	}
}

////////////////////////////////////////////////////////
// relocate to project root path
func RelocatePath(filter ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("RelocatePath failed: %w", err)
	}

	wd = strings.Replace(wd, "\\", "/", -1)
	fmt.Println("work directory: ", wd)

	var newPath string = wd

	for _, path := range filter {
		n := strings.LastIndex(wd, path)
		if n == -1 {
			continue
		}

		newPath = strings.Join([]string{wd[:n], path}, "")
		err = os.Chdir(newPath)
		if err != nil {
			return err
		}

		break
	}

	fmt.Println("relocate to new_path:", newPath)
	return nil
}

////////////////////////////////////////////////////////
// between [a, b)
func Between[T define.Number](n, a, b T) bool {
	return (n >= a && n < b)
}

////////////////////////////////////////////////////////
// Round
func Round(val float64) float64 {
	return math.Round((val*10 + 0.1) / 10)
}

////////////////////////////////////////////////////////
// pack int64
func HighId(id int64) int32 {
	return int32(id >> 32)
}

func LowId(id int64) int32 {
	return int32((id) & 0xFFFFFFFF)
}

func PackId(high, low int32) int64 {
	return (int64(high)<<32 | (int64(low) & 0x00000000FFFFFFFF))
}

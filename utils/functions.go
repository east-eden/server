package utils

import (
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog/log"
)

func CaptureException(p ...interface{}) {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		log.Error().Caller(1).Interface("exception_param", p).Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
	}
}

// relocate to project root path
func RelocatePath(filter ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("RelocatePath failed: %w", err)
	}

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

// between [a, b)
func Between(n, a, b int) bool {
	return (n >= a && n < b)
}

func BetweenInt32(n, a, b int32) bool {
	return (n >= a && n < b)
}

// Round 四舍五入
func Round(val float64) float64 {
	return math.Round((val*10 + 0.1) / 10)
}

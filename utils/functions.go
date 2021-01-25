package utils

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog/log"
)

func CaptureException() {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
	}
}

// relocate to project root path
func RelocatePath(filter ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("RelocatePath failed: %w", err)
	}

	log.Info().Str("work directory", wd).Send()

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

	log.Info().Str("new_path", newPath).Msg("relocate path success")
	return nil
}

// between [a, b)
func Between(n, a, b int) bool {
	return (n >= a && n < b)
}

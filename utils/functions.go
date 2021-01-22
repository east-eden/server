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
func RelocatePath() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("RelocatePath failed: %w", err)
	}

	log.Info().Str("work directory", wd).Send()

	var newPath string = wd
	pathFilter := []string{
		"/server",  // linux path
		"\\server", // windows path
	}

	for _, path := range pathFilter {
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

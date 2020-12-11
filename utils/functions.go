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
		log.Error().
			Interface("err", err).
			Str("stack", stack).
			Msg("panic: Recovered")
	}
}

// relocate to project root path
func RelocatePath() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("RelocatePath failed: %w", err)
	}

	var newPath string
	pathFilter := []string{
		"east-eden/server",  // linux path
		"east-eden\\server", // windows path
	}

	for _, path := range pathFilter {
		if strings.Contains(wd, path) {
			wds := strings.Split(wd, path)
			newPath = strings.Join([]string{wds[0], path}, "")
			err = os.Chdir(newPath)
			if err != nil {
				return err
			}

			break
		}
	}

	log.Info().Str("new_path", newPath).Msg("relocate path success")
	return nil
}

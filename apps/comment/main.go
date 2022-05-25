package main

import (
	"math"
	_ "net/http/pprof"
	"os"

	"github.com/east-eden/server/services/comment"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/version"
	log "github.com/rs/zerolog/log"
)

func main() {
	utils.LDFlagsCheck(os.Args, version.Version, version.Help)
	utils.Setrlimit(math.MaxUint32)

	m := comment.New()
	if err := m.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("comment run failed")
	}

	m.Stop()
}

package main

import (
	"math"
	_ "net/http/pprof"
	"os"

	"github.com/east-eden/server/services/rank"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/version"
	log "github.com/rs/zerolog/log"
)

func main() {
	utils.LDFlagsCheck(os.Args, version.Version, version.Help)
	utils.Setrlimit(math.MaxUint32)

	m := rank.New()
	if err := m.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("rank run failed")
	}

	m.Stop()
}

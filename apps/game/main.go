package main

import (
	"math"
	"os"

	"github.com/east-eden/server/services/game"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/version"
	log "github.com/rs/zerolog/log"
)

func main() {
	utils.LDFlagsCheck(os.Args, version.Version, version.Help)
	utils.Setrlimit(math.MaxUint32)

	g := game.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("game run failed")
		os.Exit(1)
	}

	g.Stop()
}

package main

import (
	"math"
	"os"

	"github.com/east-eden/server/services/combat"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/version"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.LDFlagsCheck(os.Args, version.Version, version.Help)
	utils.Setrlimit(math.MaxUint32)

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("combat run failed")
		os.Exit(1)
	}

	c.Stop()
}

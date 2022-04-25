package main

import (
	"math"
	"os"

	"github.com/east-eden/server/services/client"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/version"
	log "github.com/rs/zerolog/log"
)

var (
	BinaryVersion string
	GoVersion     string
	GitLastLog    string
)

func main() {
	utils.LDFlagsCheck(os.Args, version.Version, version.Help)
	utils.Setrlimit(math.MaxUint32)

	c := client.NewClient(nil)
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("client run error")
		os.Exit(1)
	}

	c.Stop()
}

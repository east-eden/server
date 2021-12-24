package main

import (
	"fmt"
	"math"
	"os"

	"github.com/east-eden/server/services/combat"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
)

var (
	BinaryVersion string
	GoVersion     string
	GitLastLog    string
)

func version() {
	fmt.Println("BinaryVersion:", BinaryVersion)
	fmt.Println("GoVersion:", GoVersion)
	fmt.Println("GitLastLog:", GitLastLog)
	os.Exit(0)
}

func help() {
	fmt.Println("The commands are:")
	fmt.Println("version       see all versions")
	os.Exit(0)
}

func main() {
	utils.LDFlagsCheck(os.Args, version, help)
	utils.Setrlimit(math.MaxUint32)

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("combat run failed")
		os.Exit(1)
	}

	c.Stop()
}

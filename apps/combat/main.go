package main

import (
	"fmt"
	"os"

	"bitbucket.org/funplus/server/services/combat"
	"bitbucket.org/funplus/server/utils"
	"github.com/rs/zerolog/log"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
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

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("combat run failed")
		os.Exit(1)
	}

	c.Stop()
}

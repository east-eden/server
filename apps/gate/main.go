package main

import (
	"fmt"
	"os"

	"bitbucket.org/funplus/server/services/gate"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"

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

func main() {
	utils.LDFlagsCheck(
		os.Args,

		// version
		func() {
			fmt.Println("BinaryVersion:", BinaryVersion)
			fmt.Println("GoVersion:", GoVersion)
			fmt.Println("GitLastLog:", GitLastLog)
			os.Exit(0)
		},

		// help
		func() {
			fmt.Println("The commands are:")
			fmt.Println("version       see all versions")
			os.Exit(0)
		},
	)

	g := gate.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("gate run failed")
	}

	g.Stop()
}

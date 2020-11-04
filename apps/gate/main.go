package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/gate"
	"github.com/yokaiio/yokai_server/logger"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

func main() {
	// check path
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if strings.Contains(path, "apps/") {
		os.Chdir("../../")
		newPath, _ := os.Getwd()
		fmt.Println("change current path to project root path:", newPath)
	}

	// logger init
	logger.InitLogger("gate")

	// entries init
	entries.InitEntries()

	g := gate.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("gate run failed")
	}

	g.Stop()
}

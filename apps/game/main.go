package main

import (
	"fmt"
	"os"

	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/logger"
	"bitbucket.org/east-eden/server/services/game"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

func main() {
	// relocate path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("game")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	// load xml entries
	// excel.ReadAllXmlEntries("config/entry")

	g := game.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("game run failed")
		os.Exit(1)
	}

	g.Stop()
}

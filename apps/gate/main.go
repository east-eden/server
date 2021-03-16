package main

import (
	"fmt"
	"os"

	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/services/gate"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

func main() {
	if err := utils.RelocatePath("/server", "\\server", "/server_bin", "\\server_bin"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("gate")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	// load xml entries
	// excel.ReadAllXmlEntries("config/entry")

	g := gate.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("gate run failed")
	}

	g.Stop()
}

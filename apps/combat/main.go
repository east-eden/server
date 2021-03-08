package main

import (
	"fmt"
	"os"

	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/services/combat"
	"bitbucket.org/funplus/server/utils"
	"github.com/rs/zerolog/log"

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
	logger.InitLogger("combat")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("combat run failed")
		os.Exit(1)
	}

	c.Stop()
}

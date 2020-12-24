package main

import (
	"fmt"
	"os"

	"github.com/east-eden/server/entries"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/services/gate"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

func main() {
	if err := utils.RelocatePath(); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("gate")

	// entries init
	entries.InitEntries()

	utils.GenerateFile()
	// generate excel file
	// err := excel.GenerateExcelFile("config/excel/HeroConfig.xlsx", "HeroConfig")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("generate excel file failed")
	// }

	g := gate.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("gate run failed")
	}

	g.Stop()
}

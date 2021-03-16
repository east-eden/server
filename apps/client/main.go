package main

import (
	"fmt"
	"os"

	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/services/client"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
)

func main() {
	// relocate path
	if err := utils.RelocatePath("/server", "\\server", "/server_bin", "\\server_bin"); err != nil {
		fmt.Println("relocate path failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("game")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	c := client.NewClient(nil)
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("client run error")
		os.Exit(1)
	}

	c.Stop()

}

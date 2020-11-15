package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/logger"
	"github.com/yokaiio/yokai_server/services/client"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
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
	logger.InitLogger("game")

	// entries init
	entries.InitEntries()

	c := client.NewClient(nil)
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("client run error")
		os.Exit(1)
	}

	c.Stop()

}

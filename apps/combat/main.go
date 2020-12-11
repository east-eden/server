package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/entries"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/services/combat"

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

	// relocate project path
	if strings.Contains(path, "apps/") || strings.Contains(path, "apps\\") {
		if err := os.Chdir("../../"); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		newPath, _ := os.Getwd()
		fmt.Println("change current path to project root path:", newPath)
	}

	// logger init
	logger.InitLogger("combat")

	// entries init
	entries.InitEntries()

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("combat run failed")
		os.Exit(1)
	}

	c.Stop()
}

package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/combat"
	"github.com/yokaiio/yokai_server/entries"
	_ "github.com/yokaiio/yokai_server/logger"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

func main() {
	// entries init
	entries.InitEntries()

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("combat run failed")
		os.Exit(1)
	}

	c.Stop()
}

package main

import (
	"os"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/gate"
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

	g := gate.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("gate run failed")
	}

	g.Stop()
}

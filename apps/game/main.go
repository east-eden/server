package main

import (
	"os"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

func init() {
	// set working directory as yokai_server
	os.Chdir("../../")
	logger.SetReportCaller(true)
}

func main() {
	// entries init
	entries.InitEntries()

	g := game.New()
	if err := g.Run(os.Args); err != nil {
		logger.Fatal("game run error:", err)
		os.Exit(1)
	}

	g.Stop()
}

package main

import (
	"os"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/combat"
	"github.com/yokaiio/yokai_server/entries"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

func init() {
	// set working directory as yokai_combat
	os.Chdir("../../")
	logger.SetReportCaller(true)
}

func main() {
	// entries init
	entries.InitEntries()

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		logger.Fatal("combat run error:", err)
		os.Exit(1)
	}

	c.Stop()
}

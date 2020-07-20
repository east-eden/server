package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/combat"
	"github.com/yokaiio/yokai_server/entries"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq"
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/store/consul"
	_ "github.com/micro/go-plugins/transport/grpc"
)

func init() {
	// set working directory as yokai_combat
	os.Chdir("../../")
}

func main() {
	// entries init
	entries.InitEntries()

	c := combat.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal("combat run error:", err)
		os.Exit(1)
	}

	c.Stop()
}

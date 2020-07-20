package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/gate"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq"
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/store/consul"
	_ "github.com/micro/go-plugins/transport/grpc"
)

func init() {
	// set working directory as yokai_server
	os.Chdir("../../")
}

func main() {
	// entries init
	entries.InitEntries()

	g := gate.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal("gate run error:", err)
		os.Exit(1)
	}

	g.Stop()
}

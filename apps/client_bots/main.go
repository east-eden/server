package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/client"
	"github.com/yokaiio/yokai_server/entries"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
)

func init() {
	// set working directory as yokai_server
	os.Chdir("../../")
}

func main() {
	// entries init
	entries.InitEntries()

	bots := client.NewClientBots()
	if err := bots.Run(os.Args); err != nil {
		log.Fatal("client_bots run error:", err)
		os.Exit(1)
	}

	bots.Stop()
}

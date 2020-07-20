package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/client"
	"github.com/yokaiio/yokai_server/entries"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func init() {
	// set working directory as yokai_server
	os.Chdir("../../")
}

func main() {
	// entries init
	entries.InitEntries()

	bots, err := client.NewClientBots()
	if err != nil {
		log.Fatal("client_bots new error:", err)
		os.Exit(1)
	}

	if err = bots.Run(os.Args); err != nil {
		log.Fatal("client_bots run error:", err)
		os.Exit(1)
	}

	bots.Stop()
}

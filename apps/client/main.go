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

	c := client.NewClient()
	if err := c.Run(os.Args); err != nil {
		log.Fatal("client run error:", err)
		os.Exit(1)
	}

	c.Stop()

}

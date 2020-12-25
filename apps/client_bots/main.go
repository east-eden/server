package main

import (
	"fmt"
	"log"
	"os"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/services/client"
	"github.com/east-eden/server/utils"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
)

func main() {
	// relocate path
	if err := utils.RelocatePath(); err != nil {
		fmt.Println("relocate path failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("game")

	// load excel entries
	excel.ReadAllEntries("config/excel")

	bots := client.NewClientBots()
	if err := bots.Run(os.Args); err != nil {
		log.Fatal("client_bots run error:", err)
		os.Exit(1)
	}

	bots.Stop()
}
